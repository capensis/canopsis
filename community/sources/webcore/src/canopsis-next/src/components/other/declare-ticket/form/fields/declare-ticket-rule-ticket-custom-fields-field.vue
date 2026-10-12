<template>
  <c-layout gap="3" column>
    <c-label :label="$t('common.customField', 2)" />
    <c-text-pairs-field
      :items="value"
      :text-label="$t('declareTicket.alarmFieldName')"
      :value-label="$t('declareTicket.responseField')"
      :add-button-label="$t('common.addCustomField')"
      :name="name"
      :disabled="disabled"
      text-required
      value-required
      @update:items="$emit('update:value', $event)"
    >
      <template #append-value="{ item }">
        <c-help-icon
          v-if="item.text"
          :text="$t('declareTicket.responseFieldHelpText', { field: item.text })"
          icon="help"
          color="grey darken-1"
          left
        />
      </template>
    </c-text-pairs-field>
  </c-layout>
</template>

<script>
export default {
  emits: ['update:value'],
  props: {
    value: {
      type: Array,
      required: true,
    },
    name: {
      type: String,
      default: 'declare_ticket',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
