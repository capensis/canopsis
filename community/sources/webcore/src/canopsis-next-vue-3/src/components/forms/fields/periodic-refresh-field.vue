<template>
  <v-row class="mb-3" align="start" density="compact">
    <v-col cols="5">
      <c-enabled-field
        :model-value="periodicRefresh.enabled"
        :label="label"
        hide-details
        @update:model-value="updateEnabled"
      />
    </v-col>

    <v-col cols="7">
      <c-duration-field
        :duration="periodicRefresh"
        :disabled="!periodicRefresh.enabled"
        :required="periodicRefresh.enabled"
        :name="name"
        :min="1"
        @update:duration="updateModel"
      />
    </v-col>
  </v-row>
</template>

<script>
/**
 * Field for editing view periodic refresh interval.
 *
 * @prop {Object} [periodicRefresh={}] - Periodic refresh form value.
 */
export default {
  props: {
    periodicRefresh: {
      type: Object,
      default: () => ({}),
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'periodic_refresh',
    },
  },

  emits: ['update:periodicRefresh'],

  methods: {
    updateModel(periodicRefresh) {
      this.$emit('update:periodicRefresh', periodicRefresh);
    },

    updateEnabled(enabled) {
      this.updateModel({
        ...this.periodicRefresh,
        enabled,
      });

      this.$nextTick(() => this.$validator?.validate?.(`${this.name}.value`));
    },
  },
};
</script>
