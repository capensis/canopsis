<template>
  <v-layout column>
    <v-tooltip
      v-if="!hideAckResources"
      top
    >
      <template #activator="{ on }">
        <v-checkbox
          v-on="on"
          v-field="form.ack_resources"
          :label="$t('alarm.ackResources')"
          color="primary"
        />
      </template>
      <span>{{ $t('alarm.ackResourcesQuestion') }}</span>
    </v-tooltip>
    <c-description-field
      v-field="form.comment"
      :label="$t('common.note')"
      :required="isNoteRequired"
      name="comment"
    />
  </v-layout>
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isNoteRequired: {
      type: Boolean,
      default: false,
    },
    hideAckResources: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
