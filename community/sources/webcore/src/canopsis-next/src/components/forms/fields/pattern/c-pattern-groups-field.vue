<template lang="pug">
  v-layout(column)
    v-flex(xs12)
      v-alert(
        v-if="!groups.length",
        :value="true",
        type="info"
      ) {{ disabled || readonly ? $t('pattern.noDataDisabled') : $t('pattern.noData') }}
    v-layout(v-for="(group, index) in groups", :key="group.key", wrap, row)
      v-flex(xs12)
        c-pattern-group-field(
          v-field="groups[index]",
          :attributes="attributes",
          :disabled="disabled",
          :readonly="readonly",
          @remove="removeItemFromArray(index)"
        )
      v-layout(v-show="index !== groups.length - 1", justify-center)
        c-pattern-operator-chip {{ $t('common.or') }}
    v-layout(v-if="!readonly", row, align-center)
      v-btn.ml-0(
        :color="hasGroupsErrors ? 'error' : 'primary'",
        :disabled="disabled",
        @click="addFilterGroup"
      ) {{ $t('pattern.addGroup') }}
      span.error--text(v-show="hasGroupsErrors") {{ $t('pattern.errors.groupRequired') }}
</template>

<script>
import { patternRulesToGroup } from '@/helpers/forms/pattern';

import { formArrayMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin],
  model: {
    prop: 'groups',
    event: 'input',
  },
  props: {
    groups: {
      type: Array,
      required: true,
    },
    attributes: {
      type: Array,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'groups',
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    hasGroupsErrors() {
      return this.errors.has(this.name);
    },
  },
  watch: {
    groups() {
      this.$nextTick(() => {
        if (this.required) {
          this.$validator.validate(this.name);
        }
      });
    },
    required: {
      immediate: true,
      handler(value) {
        if (value) {
          this.attachMinValueRule();
        } else {
          this.detachMinValueRule();
        }
      },
    },
  },
  beforeDestroy() {
    this.detachMinValueRule();
  },
  methods: {
    attachMinValueRule() {
      this.$validator.attach({
        name: this.name,
        rules: 'min_value:1',
        getter: () => this.groups.length,
        vm: this,
      });
    },

    detachMinValueRule() {
      this.$validator.detach(this.name);
    },

    addFilterGroup() {
      const [firstAttribute] = this.attributes;

      this.addItemIntoArray(patternRulesToGroup([
        { field: firstAttribute?.value },
      ]));
    },
  },
};
</script>
