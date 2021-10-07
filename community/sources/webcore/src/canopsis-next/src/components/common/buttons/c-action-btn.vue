<template lang="pug">
  v-tooltip(bottom)
    slot(slot="activator", name="button")
      v-btn.mx-1(
        :disabled="disabled",
        :data-test="$attrs['data-test']",
        :loading="loading",
        icon,
        @click.stop.prevent="$listeners.click"
      )
        v-icon(:color="preparedProps.color") {{ preparedProps.icon }}
    span {{ preparedProps.tooltip }}
</template>

<script>
export default {
  props: {
    type: {
      type: String,
      default: null,
      validator: value => ['edit', 'duplicate', 'delete'].includes(value),
    },
    icon: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: '',
    },
    tooltip: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    preparedProps() {
      const typesPresetsMap = {
        edit: {
          tooltip: this.$t('common.edit'),
          icon: 'edit',
        },
        duplicate: {
          tooltip: this.$t('common.duplicate'),
          icon: 'file_copy',
        },
        delete: {
          tooltip: this.$t('common.delete'),
          icon: 'delete',
          color: 'error',
        },
      };

      const props = typesPresetsMap[this.type] || {};

      return {
        icon: this.icon || props.icon,
        color: this.color || props.color,
        tooltip: this.tooltip || props.tooltip,
      };
    },
  },
};
</script>
