<template lang="pug">
  v-card.my-1.pa-2
    v-radio-group(
    :input-value="condition",
    @change="$emit('update:condition', $event)",
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
        v-if="!initialGroup",
        @click="$emit('deleteGroup')",
        color="red darken-4",
        outline,
        block,
        small,
        flat
        ) {{$t("filterEditor.buttons.deleteGroup")}}

    div(v-for="(rule, ruleKey) in rules", :key="ruleKey")
      filter-rule(
      @deleteRule="deleteRule(ruleKey)",
      @update:rule="updateRule(ruleKey, $event)",
      :rule="rule",
      :operators="operators",
      :possibleFields="possibleFields",
      )

    div(v-for="(group, groupKey) in groups", :key="groupKey")
      filter-group.filterGroup(
      @deleteGroup="deleteGroup(groupKey)",
      :condition.sync="group.condition",
      :rules.sync="group.rules",
      :groups.sync="group.groups",
      :possibleFields="possibleFields",
      )
</template>

<script>
import omit from 'lodash/omit';
import cloneDeep from 'lodash/cloneDeep';

import { FILTER_OPERATORS, FILTER_DEFAULT_VALUES } from '@/constants';
import uid from '@/helpers/uid';

import FilterRule from './filter-rule.vue';

/**
 * Component representing a group in MongoDB filter
 *
 * @prop {Array} possibleFields - Boolean to determine if it's the root filter's group
 * @prop {boolean} [initialGroup=false] - Boolean to determine if it's the root filter's group
 * @prop {string} [condition='$and'] - Base condition of the group : "$and" or "$or"
 * @prop {Array} [rules=[]] - Rules of the current group
 * @prop {Array} [groups=[]] - Groups of the current group
 *
 * @event condition#update
 * @event deleteGroup#click
 */
export default {
  name: 'filter-group', // We need it for recursive
  components: {
    FilterRule,
  },
  props: {
    possibleFields: {
      type: Array,
      required: true,
    },
    initialGroup: {
      type: Boolean,
      default: false,
    },
    condition: {
      type: String,
      default: FILTER_DEFAULT_VALUES.condition,
    },
    rules: {
      type: Object,
      default() {
        return {};
      },
    },
    groups: {
      type: Object,
      default() {
        return {};
      },
    },
  },
  data() {
    return {
      operators: Object.values(FILTER_OPERATORS),
    };
  },
  methods: {
    updateRule(key, value) {
      this.$emit('update:rules', { ...this.rules, [key]: value });
    },

    /**
     * @description Invoked on a click on 'Add Rule' button. Add an empty object to the 'rules' array
     */
    addRule() {
      this.$emit('update:rules', {
        ...this.rules,
        [uid('rule')]: cloneDeep(FILTER_DEFAULT_VALUES.rule),
      });
    },

    /**
     * @description Invoked on a click on 'Add Group' button. Add a Group to the 'groups' array
     */
    addGroup() {
      this.$emit('update:groups', {
        ...this.groups,
        [uid('group')]: cloneDeep(FILTER_DEFAULT_VALUES.group),
      });
    },

    /**
     * @description Invoked when a 'deleteRule' event is fired. Delete a rule from the 'rules' array
     * @param {string} key
     */
    deleteRule(key) {
      this.$emit('update:rules', omit(this.rules, [key]));
    },

    /**
     * @description Invoked when a 'deleteGroup' event is fired. Delete a group from the 'groups' array
     * @param {string} key
     */
    deleteGroup(key) {
      this.$emit('update:groups', omit(this.groups, [key]));
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
