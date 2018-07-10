<template lang="pug">
    v-container.filterGroup.pa-2(
      fluid,
    )
      v-radio-group(
        mandatory,
        row,
        hide-details,
        :input-value="condition",
        @change="$emit('update:condition', $event)"
      )
        v-radio(label="AND", value="$and", color="blue darken-4")
        v-radio(label="OR", value="$or", color="blue darken-4")

      v-layout.text-xs-center(wrap, justify-space-around)
        v-flex(xs5, md3)
          v-btn(
            @click="handleAddRuleClick",
            flat,
            outline,
            block,
            small
          ) {{$t("mFilterEditor.buttons.addRule")}}
        v-flex(xs5, md3)
          v-btn(
            @click="handleAddGroupClick",
            flat,
            outline,
            block,
            small
          ) {{$t("mFilterEditor.buttons.addGroup")}}
        v-flex(xs5, md3)
          v-btn(
            v-if="!initialGroup",
            @click="handleDeleteGroupClick",
            flat,
            outline,
            color="red darken-4",
            block,
            small
          ) {{$t("mFilterEditor.buttons.deleteGroup")}}

      div(
        v-for="(rule, index) in rules",
        :key="'rule-' + index",
      )
        filter-rule(
          @deleteRuleClick="handleDeleteRuleClick",
          :index="index",
          :field.sync="rule.field",
          :operator.sync="rule.operator",
          :input.sync="rule.input",
          :operators="operators",
          :possibleFields="possibleFields",
        )

      div(
        v-for="(group, index) in groups",
        :key="'group-' + index"
      )
        filter-group.filterGroup(
          @deleteGroup="deleteGroup",
          :index="index",
          :condition.sync="group.condition",
          :possibleFields="possibleFields",
          :operators="operators",
          :rules="group.rules",
          :groups="group.groups",
        )
</template>

<script>
import FilterRule from './filter-rule.vue';

/**
 * Component representing a group in MongoDB filter
 *
 * @prop {Boolean} [initialGroup] - Boolean to determine if it's the root filter's group
 * @prop {Number} [index] - Index of the group
 * @prop {String} [condition] - Base condition of the group : "$and" or "$or"
 *
 * @event condition#update
 * @event deleteGroup#click
 */
export default {
  name: 'filter-group',
  components: {
    FilterRule,
  },
  props: {
    initialGroup: {
      type: Boolean,
    },
    index: {
      type: Number,
      required: true,
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
      operators: [
        { value: 'equal' },
        { value: 'not equal' },
        { value: 'in' },
        { value: 'not in' },
        { value: 'begins with' },
        { value: 'doesn\'t begin with' },
        { value: 'contains' },
        { value: 'doesn\'t contain' },
        { value: 'ends with' },
        { value: 'doesn\'t end with' },
        { value: 'is empty' },
        { value: 'is not empty' },
        { value: 'is null' },
        { value: 'is not null' },
      ],
    };
  },
  methods: {
    /**
     * @description Invoked on a click on 'Add Rule' button. Add an empty object to the 'rules' array
     */
    handleAddRuleClick() {
      this.rules.push({
        field: '',
        operator: '',
        input: '',
      });
    },
    /**
     * @description Invoked when a 'deleteRuleClick' event is fired. Delete a rule from the 'rules' array
     * @param Int
     */
    handleDeleteRuleClick(index) {
      this.rules.splice(index, 1);
    },
    /**
     * @description Invoked on a click on 'Add Group' button. Add a Group to the 'groups' array
     */
    handleAddGroupClick() {
      this.groups.push({
        condition: '$and',
        groups: [],
        rules: [],
      });
    },
    /**
     * @description Invoked on a click on 'Delete this group' button. Emit a 'deleteGroup'
     * event to the parent, that will actually delete the group
     */
    handleDeleteGroupClick() {
      this.$emit('deleteGroup', this.index);
    },
    /**
     * @description Invoked when a 'deleteGroup' event is fired. Delete a group from the 'groups' array
     * @param INT
     */
    deleteGroup(index) {
      this.groups.splice(index, 1);
    },
  },
};
</script>

<style scoped lang="scss">
  .filterGroup {
    border: 1px solid gray;
  }

  button {
    @media (max-width:'500px') {
      font-size: 0.6em;
    }
  }
</style>
