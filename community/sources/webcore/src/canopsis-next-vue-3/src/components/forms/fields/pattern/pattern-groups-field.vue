<template>
  <c-layout gap="3" column>
    <v-alert
      v-if="!groups.length"
      type="info"
    >
      {{ disabled || readonly ? $t('pattern.noDataDisabled') : $t('pattern.noData') }}
    </v-alert>
    <c-layout
      v-for="(group, index) in groups"
      :key="group.key"
      wrap
    >
      <pattern-group-field
        :group="groups[index]"
        :attributes="attributes"
        :disabled="disabled"
        :readonly="readonly"
        @update:group="updateItemInArray(index, $event)"
        @remove="removeItemFromArray(index)"
      />
      <c-layout
        v-show="index !== groups.length - 1"
        justify-center
      >
        <c-pattern-operator-chip>{{ $t('common.or') }}</c-pattern-operator-chip>
      </c-layout>
    </c-layout>

    <div class="align-self-start">
      <c-btn-with-error
        v-if="!readonly"
        :error="hasGroupsErrors ? $t('pattern.errors.groupRequired'): ''"
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
import { computed } from 'vue';
import { useField } from 'vee-validate';

import { patternRulesToGroup } from '@/helpers/entities/pattern/form';

import { formArrayMixin } from '@/mixins/form';

import PatternGroupField from './pattern-group-field.vue';

export default {
  components: { PatternGroupField },
  mixins: [formArrayMixin],
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
  setup(props) {
    const fieldName = computed(() => props.name);
    const groupsLength = computed(() => (props.groups ?? []).length);
    const {
      errorMessage,
      validate,
    } = useField(
      fieldName,
      () => (!props.required || groupsLength.value ? true : 'Required'),
      {
        initialValue: groupsLength.value,
        validateOnValueUpdate: false,
      },
    );

    return {
      groupErrorMessage: errorMessage,
      validateGroupsField: validate,
    };
  },
  computed: {
    hasGroupsErrors() {
      return !!this.groupErrorMessage;
    },
  },
  watch: {
    groups() {
      this.$nextTick(() => {
        if (this.required) {
          this.validateGroupsField();
        }
      });
    },
  },
  methods: {
    addFilterGroup() {
      const [firstAttribute] = this.attributes;

      this.addItemIntoArray(patternRulesToGroup([
        { field: firstAttribute?.value },
      ]));
    },
  },
};
</script>
