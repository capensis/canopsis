<template>
  <v-layout class="gap-3" column>
    <c-label>{{ $tc('common.customField', 2) }}</c-label>
    <c-text-pairs-field
      v-field="value"
      :text-label="$t('declareTicket.alarmFieldName')"
      :value-label="$t('declareTicket.responseField')"
      :name="name"
      :disabled="disabled"
      text-required
      value-required
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
  </v-layout>
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
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
