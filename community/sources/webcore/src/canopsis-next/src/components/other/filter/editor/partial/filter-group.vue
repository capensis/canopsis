<template lang="pug">
  v-card.my-1.pa-2(data-test="filterGroup")
    c-operator-field(v-field="group.condition")
    v-layout.text-xs-center(wrap, justify-space-around)
      v-flex(xs5, md3)
        v-btn(
          data-test="addRule",
          outline,
          block,
          small,
          flat,
          @click="addRule"
        ) {{ $t("filterEditor.buttons.addRule") }}
      v-flex(xs5, md3)
        v-btn(
          data-test="addGroup",
          outline,
          block,
          small,
          flat,
          @click="addGroup"
        ) {{ $t("filterEditor.buttons.addGroup") }}
      v-flex(xs5, md3)
        v-btn(
          data-test="deleteGroup",
          v-if="!isInitial",
          color="red darken-4",
          outline,
          block,
          small,
          flat,
          @click="$emit('delete-group')"
        ) {{ $t("filterEditor.buttons.deleteGroup") }}

    v-layout
      pattern-information.my-2.mr-2(v-if="filterInformationShown") {{ groupOperator }}
      v-flex
        div(data-test="filterRuleLayout")
          filter-rule(
            v-for="(rule, ruleKey) in group.rules",
            :key="ruleKey",
            :rule="rule",
            :name="ruleKey",
            :possible-fields="possibleFields",
            @delete-rule="deleteRule(ruleKey)",
            @input="updateRule(ruleKey, $event)"
          )

        div(data-test="filterGroupLayout")
          filter-group.filterGroup(
            v-for="(group, groupKey) in group.groups",
            :key="groupKey",
            :group="group",
            :possible-fields="possibleFields",
            @input="updateGroup(groupKey, $event)",
            @delete-group="deleteGroup(groupKey)"
          )
</template>

<script>
import { omit, cloneDeep } from 'lodash';

import { FILTER_DEFAULT_VALUES, FILTER_MONGO_OPERATORS } from '@/constants';

import uid from '@/helpers/uid';

import { formMixin } from '@/mixins/form';

import PatternInformation from '@/components/other/pattern/pattern-information.vue';

import FilterRule from './filter-rule.vue';

/**
 * Component representing a group in MongoDB filter
 */
export default {
  name: 'FilterGroup', // We need it for recursive
  components: {
    PatternInformation,
    FilterRule,
  },
  mixins: [formMixin],
  model: {
    prop: 'group',
    event: 'input',
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
  computed: {
    filterInformationShown() {
      return Object.keys(this.group.rules).length > 1;
    },

    groupOperator() {
      return this.group.condition === FILTER_MONGO_OPERATORS.and ? this.$t('common.and') : this.$t('common.or');
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
