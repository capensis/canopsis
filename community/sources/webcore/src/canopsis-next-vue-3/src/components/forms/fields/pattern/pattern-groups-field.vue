<template>
  <c-layout gap="3" column>
    <v-alert v-if="!groups.length" type="info">
      {{ disabled || readonly ? $t('pattern.noDataDisabled') : $t('pattern.noData') }}
    </v-alert>
    <c-layout v-for="(group, index) in groups" :key="group.key" wrap>
      <pattern-group-field
        :group="groups[index]"
        :attributes="attributes"
        :disabled="disabled"
        :readonly="readonly"
        @update:group="updateItemInArray(index, $event)"
        @remove="removeItemFromArray(index)"
      />
      <c-layout v-show="index !== groups.length - 1" justify-center>
        <c-pattern-operator-chip>{{ $t('common.or') }}</c-pattern-operator-chip>
      </c-layout>
    </c-layout>

    <div class="align-self-start">
      <c-btn-with-error
        v-if="!readonly"
        :error="hasGroupsErrors ? $t('pattern.errors.groupRequired') : ''"
        :disabled="disabled"
        class="mt-2"
        @click="addFilterGroup"
      >
        {{ $t('pattern.addGroup') }}
      </c-btn-with-error>
    </div>
  </c-layout>
</template>

<script>
import { patternRulesToGroup } from '@/helpers/entities/pattern/form';

import { formArrayMixin } from '@/mixins/form';
import { validatorMixin } from '@/mixins/validator';

import PatternGroupField from './pattern-group-field.vue';

export default {
  components: { PatternGroupField },
  mixins: [formArrayMixin, validatorMixin],
  emits: ['update:groups'],
  props: {
    groups: {
      type: Array,
      required: true,
    },
    attributes: {
      type: Array,
      default: () => [],
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
        if (!this.required) {
          this.clearGroupsErrors();

          return;
        }

        if (this.groups.length) {
          this.clearGroupsErrors();
        } else {
          this.validateGroupsField();
        }
      });
    },
  },
  methods: {
    clearGroupsErrors() {
      this.errors.clear(this.name);
    },

    validateGroupsField() {
      if (!this.required || this.groups.length) {
        this.clearGroupsErrors();

        return true;
      }

      this.errors.add({
        field: this.name,
        msg: this.$t('pattern.errors.groupRequired'),
      });

      return false;
    },

    addFilterGroup() {
      const [firstAttribute] = this.attributes;

      this.clearGroupsErrors();
      this.addItemIntoArray(patternRulesToGroup([{ field: firstAttribute?.value }]));
    },
  },
};
</script>
