<template lang="pug">
  v-card.my-1.pa-2()
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
        ) {{$t("mFilterEditor.buttons.addRule")}}
      v-flex(xs5, md3)
        v-btn(
        @click="addGroup",
        outline,
        block,
        small,
        flat
        ) {{$t("mFilterEditor.buttons.addGroup")}}
      v-flex(xs5, md3)
        v-btn(
        v-if="!initialGroup",
        @click="$emit('deleteGroup', index)",
        color="red darken-4",
        outline,
        block,
        small,
        flat
        ) {{$t("mFilterEditor.buttons.deleteGroup")}}

    div(v-for="(rule, index) in rules", :key="`rule-${index}`")
      filter-rule(
      @deleteRule="deleteRule",
      :index="index",
      :field.sync="rule.field",
      :operator.sync="rule.operator",
      :input.sync="rule.input",
      :operators="operators",
      :possibleFields="possibleFields",
      )

    div(v-for="(group, index) in groups", :key="`group-${index}`")
      filter-group.filterGroup(
      @deleteGroup="deleteGroup",
      :index="index",
      :condition.sync="group.condition",
      :possibleFields="possibleFields",
      :rules="group.rules",
      :groups="group.groups",
      )
</template>

<script>
import { OPERATORS } from '@/constants';

import FilterRule from './filter-rule.vue';

/**
 * Component representing a group in MongoDB filter
 *
 * @prop {Boolean} [initialGroup] - Boolean to determine if it's the root filter's group
 * @prop {Number} index - Index of the group
 * @prop {String} [condition] - Base condition of the group : "$and" or "$or"
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
    initialGroup: {
      type: Boolean,
    },
    index: {
      type: Number,
      default: 0,
    },
    condition: {
      type: String,
      default() {
        return '$and';
      },
    },
    rules: {
      type: Array,
      default() {
        return [];
      },
    },
    groups: {
      type: Array,
      default() {
        return [];
      },
    },
    possibleFields: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      operators: Object.values(OPERATORS),
    };
  },
  methods: {
    /**
     * @description Invoked on a click on 'Add Rule' button. Add an empty object to the 'rules' array
     */
    addRule() {
      this.rules.push({
        field: '',
        operator: '',
        input: '',
      });
    },

    /**
     * @description Invoked on a click on 'Add Group' button. Add a Group to the 'groups' array
     */
    addGroup() {
      this.groups.push({
        condition: '$and',
        groups: [],
        rules: [],
      });
    },

    /**
     * @description Invoked when a 'deleteRule' event is fired. Delete a rule from the 'rules' array
     * @param {number} index
     */
    deleteRule(index) {
      this.rules.splice(index, 1);
    },

    /**
     * @description Invoked when a 'deleteGroup' event is fired. Delete a group from the 'groups' array
     * @param {number} index
     */
    deleteGroup(index) {
      this.groups.splice(index, 1);
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
