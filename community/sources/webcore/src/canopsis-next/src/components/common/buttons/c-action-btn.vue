<template>
  <c-simple-tooltip
    :top="top"
    :right="right"
    :bottom="bottom"
    :left="left"
    :disabled="disabled"
    :content="preparedProps.tooltip"
  >
    <template #activator="{ on: tooltipOn }">
      <slot
        :on="tooltipOn"
        name="button"
      >
        <v-btn
          :disabled="disabled"
          :loading="loading"
          :small="small"
          :color="btnColor"
          :dark="dark"
          :input-value="inputValue"
          class="mx-1 my-0 c-action-btn__button"
          icon
          v-on="tooltipOn"
          @click.stop.prevent="$listeners.click"
        >
          <v-icon :color="preparedProps.color">
            {{ preparedProps.icon }}
          </v-icon>
        </v-btn>
      </slot>
    </template>
  </c-simple-tooltip>
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
    top: {
      type: Boolean,
      required: false,
    },
    right: {
      type: Boolean,
      required: false,
    },
    left: {
      type: Boolean,
      required: false,
    },
    bottom: {
      type: Boolean,
      required: false,
      default() {
        return !this.top && !this.right && !this.left;
      },
    },
    small: {
      type: Boolean,
      required: false,
    },
    btnColor: {
      type: String,
      required: false,
    },
    dark: {
      type: Boolean,
      default: false,
    },
    inputValue: {
      type: Boolean,
      required: false,
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
