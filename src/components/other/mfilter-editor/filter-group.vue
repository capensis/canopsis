<template lang="pug">
    v-container(
      fluid,
      class="filterGroup pa-2"
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

      v-layout(wrap, class="text-xs-center", justify-space-around)
        v-flex(xs5, md3)
          v-btn(
            @click="handleAddRuleClick",
            flat,
            outline,
            block,
            small
          ) {{$t("m_filter_editor.buttons.add_rule")}}
        v-flex(xs5, md3)
          v-btn(
            @click="handleAddGroupClick",
            flat,
            outline,
            block,
            small
          ) {{$t("m_filter_editor.buttons.add_group")}}
        v-flex(xs5, md3)
          v-btn(
            v-if="!initialGroup",
            @click="handleDeleteGroupClick",
            flat,
            outline,
            color="red darken-4",
            block,
            small
          ) {{$t("m_filter_editor.buttons.delete_group")}}

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
        filter-group(
          class="filterGroup",
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

export default {
  name: 'filter-group',
  components: {
    FilterRule,
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
