<template>
  <v-list-group>
    <template #activator="">
      <v-list-item>
        {{ label }}
        <div class="font-italic text-caption ml-1">
          ({{ $t('common.optional') }})
        </div>
      </v-list-item>
    </template>
    <v-container>
      <c-enabled-field v-field="value.enabled" />
      <v-text-field
        v-field="value.value"
        :label="$t('common.output')"
        :disabled="!value.enabled"
      />
    </v-container>
  </v-list-group>
</template>

<script>
/**
 * Component to set fast-ack output
 *
 * @prop {Object} [value] - Value
 */
export default {
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    label: {
      type: String,
      default: '',
    },
  },
};
</script>
