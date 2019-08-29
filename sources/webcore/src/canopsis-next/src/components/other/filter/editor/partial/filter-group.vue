<template lang="pug">
  v-card.my-1.pa-2
    v-radio-group(
      :value="group.condition",
      @change="updateField('condition', $event)",
      hide-details,
      mandatory,
      row
    )
      v-radio(label="AND", value="$and", color="blue darken-4")
      v-radio(label="OR", value="$or", color="blue darken-4")
    v-layout.text-xs-center(wrap, justify-space-around)
      v-flex(xs5, md3)
        v-btn(
          @click="addRule",
          outline,
          block,
          small,
          flat
        ) {{$t("filterEditor.buttons.addRule")}}
      v-flex(xs5, md3)
        v-btn(
          @click="addGroup",
          outline,
          block,
          small,
          flat
        ) {{$t("filterEditor.buttons.addGroup")}}
      v-flex(xs5, md3)
        v-btn(
          v-if="!isInitial",
          @click="$emit('deleteGroup')",
          color="red darken-4",
          outline,
          block,
          small,
          flat
        ) {{$t("filterEditor.buttons.deleteGroup")}}

    div(v-for="(rule, ruleKey) in group.rules", :key="ruleKey")
      filter-rule(
        :rule="rule",
        :possibleFields="possibleFields",
        @deleteRule="deleteRule(ruleKey)",
        @update:rule="updateRule(ruleKey, $event)"
      )

    div(v-for="(group, groupKey) in group.groups", :key="groupKey")
      filter-group.filterGroup(
        :group="group",
        :possibleFields="possibleFields",
        @deleteGroup="deleteGroup(groupKey)",
        @update:group="updateGroup(groupKey, $event)"
      )
</template>

<script>
import { omit, cloneDeep } from 'lodash';

import { FILTER_DEFAULT_VALUES } from '@/constants';

import uid from '@/helpers/uid';
import formMixin from '@/mixins/form';

import FilterRule from './filter-rule.vue';

/**
 * Component representing a group in MongoDB filter
 *
 * @prop {Array} possibleFields - Boolean to determine if it's the root filter's group
 * @prop {Object} group - Group object
 * @prop {boolean} [isInitial=false] - Boolean to determine if it's the root filter's group
 *
 * @event group#update
 * @event deleteGroup#click
 */
export default {
  name: 'filter-group', // We need it for recursive
  components: {
    FilterRule,
  },
  mixins: [formMixin],
  model: {
    prop: 'group',
    event: 'update:group',
  },
  props: {
    possibleFields: {
      type: Array,
      required: true,
    },
    group: {
      type: Object,
      required: true,
      validator: value => value.condition && value.groups && value.rules,
    },
    isInitial: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    updateRule(key, value) {
      this.updateField('rules', { ...this.group.rules, [key]: value });
    },

    updateGroup(key, value) {
      this.updateField('groups', { ...this.group.groups, [key]: value });
    },

    /**
     * @description Invoked on a click on 'Add Rule' button. Add an empty object to the 'rules' array
     */
    addRule() {
      this.updateRule(uid('rule'), cloneDeep(FILTER_DEFAULT_VALUES.rule));
    },

    /**
     * @description Invoked on a click on 'Add Group' button. Add a Group to the 'groups' array
     */
    addGroup() {
      this.updateGroup(uid('group'), cloneDeep(FILTER_DEFAULT_VALUES.group));
    },

    /**
     * @description Invoked when a 'deleteRule' event is fired. Delete a rule from the 'rules' array
     * @param {string} key
     */
    deleteRule(key) {
      this.updateField('rules', omit(this.group.rules, [key]));
    },

    /**
     * @description Invoked when a 'deleteGroup' event is fired. Delete a group from the 'groups' array
     * @param {string} key
     */
    deleteGroup(key) {
      this.updateField('groups', omit(this.group.groups, [key]));
    },
  },
};
</script>

<style scoped lang="scss">
  button {
    @media (max-width: 500px) {
      font-size: 0.6em;
    }
  }
</style>
