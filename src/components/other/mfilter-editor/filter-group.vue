<template lang='pug'>
    v-container(
      class='filterGroup'
    )
      v-radio-group(
        :input-value='condition'
        @change="$emit('update:condition', $event)"
      )
        v-radio(label='AND' value='$and')
        v-radio(label='OR' value='$or')

      v-btn(@click='handleAddRuleClick') {{$t('m_filter_editor.buttons.add_rule')}}
      v-btn(@click='handleAddGroupClick') {{$t('m_filter_editor.buttons.add_group')}}

      template(v-if='!initialGroup')
        v-btn(@click='handleDeleteGroupClick') {{$t('m_filter_editor.buttons.delete_group')}}

      div(
        v-for='(group, index) in groups'
        :key="'group-' + index"
      )
        filter-group(
          @deleteGroup='deleteGroup'
          :index='index'
          :condition.sync='group.condition'
          :possibleFields='possibleFields'
          :operators='operators'
          :rules='group.rules'
          :groups='group.groups'
        )

      div(
        v-for='(rule, index) in rules'
        :key="'rule-' + index"
      )
        filter-rule(
          @deleteRuleClick='handleDeleteRuleClick'
          :index='index'
          :field.sync='rule.field'
          :operator.sync='rule.operator'
          :input.sync='rule.input'
          :isValid.sync='rule.isValid'
          :operators='operators'
          :possibleFields='possibleFields'
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
      required: true,
    },
    rules: {
      type: Array,
      required: true,
    },
    groups: {
      type: Array,
      required: true,
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

<style scoped>
    .filterGroup {
        min-width: 500px;
        border: 1px solid lightgray;
    }
</style>
