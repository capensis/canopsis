<template lang="pug">
  v-tooltip(
    :top="top",
    :right="right",
    :bottom="bottom",
    :left="left",
    :disabled="disabled",
    :custom-activator="customTooltipActivator",
    :lazy="lazy"
  )
    template(#activator="{ on: tooltipOn }")
      slot(name="button", :on="tooltipOn")
        v-badge.c-action-btn__badge(v-if="badgeValue", :color="badgeColor", overlap)
          template(#badge="")
            v-tooltip(
              :top="top",
              :right="right",
              :bottom="bottom",
              :left="left",
              :disabled="!badgeTooltip",
              :custom-activator="customTooltipActivator"
            )
              template(#activator="{ on: badgeTooltipOn }")
                slot(name="badgeIcon", :on="badgeTooltipOn")
                  v-icon(v-on="badgeTooltipOn", color="white") {{ badgeIcon }}
              span {{ badgeTooltip }}
          v-btn.ma-0.c-action-btn__button(
            v-on="tooltipOn",
            :disabled="disabled",
            :loading="loading",
            :small="small",
            :color="btnColor",
            icon,
            @click.stop.prevent="$listeners.click"
          )
            v-icon(:color="preparedProps.color") {{ preparedProps.icon }}
        v-btn.mx-1.my-0.c-action-btn__button(
          v-else,
          v-on="tooltipOn",
          :disabled="disabled",
          :loading="loading",
          :small="small",
          :color="btnColor",
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
    badgeValue: {
      type: Boolean,
      default: false,
    },
    badgeIcon: {
      type: String,
      default: 'priority_high',
    },
    badgeColor: {
      type: String,
      default: 'error',
    },
    badgeTooltip: {
      type: String,
      default: '',
    },
    customTooltipActivator: {
      type: Boolean,
      default: false,
    },
    lazy: {
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

<style lang="scss" scoped>
.c-action-btn__badge {
  margin: 6px 4px;

  & ::v-deep .v-badge__badge {
    font-size: 11px;
    top: -4px;
    right: -4px;
    height: 16px;
    width: 16px;

    .v-icon {
      font-size: 11px;
    }
  }
}
</style>
