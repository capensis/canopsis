<template>
  <v-card :class="wrapperClasses" class="modal-wrapper position-relative">
    <modal-mass-actions-panel :id="modal.id" />

    <v-card-title
      v-if="$slots.title"
      :style="titleStyle"
      :class="titleClass ?? 'text-white'"
      class="modal-wrapper__title-container"
    >
      <div class="modal-wrapper__title text-h5 ga-1">
        <div>
          <slot name="title" />
        </div>

        <div :class="titleButtonsClasses" class="modal-title-buttons">
          <div v-if="minimize" class="modal-title-button__wrapper">
            <v-tooltip v-if="!modal.minimized" location="left">
              <template #activator="{ props: activatorProps }">
                <v-btn
                  class="v-btn--minimize my-0"
                  size="default"
                  variant="text"
                  icon
                  v-bind="activatorProps"
                  @click="minimizeModal"
                >
                  <c-icon size="large">
minimize
</c-icon>
                </v-btn>
              </template>

              <span>{{ t('modals.common.titleButtons.minimizeTooltip') }}</span>
            </v-tooltip>

            <v-btn
v-else
size="x-small"
variant="text"
icon
@click="maximizeModal">
              <c-icon> maximize </c-icon>
            </v-btn>
          </div>

          <div v-if="close" class="modal-title-button__wrapper">
            <v-btn
              :size="modal.minimized ? 'x-small' : 'default'"
              variant="text"
              icon
              @click="closeModal"
            >
              <c-icon :size="modal.minimized ? undefined : 'large'">
close
</c-icon>
            </v-btn>
          </div>
        </div>
      </div>
    </v-card-title>

    <template v-if="!modal.minimized">
      <v-card-text v-if="$slots.text" key="text" :class="textClass">
        <slot name="text" />
      </v-card-text>

      <template v-if="$slots.actions">
        <c-layout column>
          <v-divider key="divider" />

          <v-card-actions key="actions" class="justify-end align-center pa-3">
            <v-defaults-provider :defaults="actionButtonDefaults">
              <slot name="actions" />
            </v-defaults-provider>
          </v-card-actions>
        </c-layout>
      </template>
    </template>
  </v-card>
</template>

<script>
import { computed, inject } from 'vue';

import { useModals } from '@/hooks/modals';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    fillHeight: {
      type: Boolean,
      default: false,
    },
    minimize: {
      type: Boolean,
      default: false,
    },
    close: {
      type: [Boolean, Function],
      default: false,
    },
    titleColor: {
      type: String,
      default: 'rgb(var(--v-theme-primary))',
    },
    titleClass: {
      type: String,
      required: false,
    },
    textClass: {
      type: String,
      required: false,
    },
  },

  setup(props) {
    const modal = inject('$modal');
    const injectedCloseModal = inject('$closeModal', null);
    const { dialogPropsMap = {} } = inject('$modalsOptions', {});
    const modals = useModals();
    const { t } = useI18n();

    const isAutoHeight = computed(() => Boolean(dialogPropsMap[modal.name]?.autoHeight));

    const titleStyle = computed(() => ({
      backgroundColor: props.titleColor,
    }));

    const wrapperClasses = computed(() => ({
      'fill-min-height': props.fillHeight,
      'modal-wrapper--auto-height': isAutoHeight.value,
    }));

    const actionButtonDefaults = {
      VBtn: {
        variant: 'elevated',
      },
    };

    const closeModal = () => {
      if (props.close !== true) {
        return props.close();
      }

      return injectedCloseModal ? injectedCloseModal() : modals.hide({ id: modal.id });
    };

    const minimizeModal = () => modals.minimize({ id: modal.id });

    const maximizeModal = () => modals.maximize({ id: modal.id });

    const titleButtonsClasses = computed(() => ({
      close: props.close,
      minimize: props.minimize,
    }));

    return {
      actionButtonDefaults,
      modal,
      t,
      titleStyle,
      wrapperClasses,
      closeModal,
      minimizeModal,
      maximizeModal,
      titleButtonsClasses,
    };
  },
};
</script>

<style lang="scss" scoped>
.modal-wrapper {
  position: relative;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;

  &--auto-height {
    height: auto;
    max-height: 100%;
    min-height: 0;
    overflow: hidden;

    .v-card-text {
      flex: 1 1 auto;
      min-height: 0;
      overflow-y: auto;
    }
  }

  &__title {
    display: flex;
    justify-content: space-between;
    width: 100%;
    align-items: center;

    &-container {
      padding: 16px 24px 10px;
    }

    & > div {
      display: flex;
      align-items: center;
    }
  }

  .modal-title-buttons {
    display: flex;
    align-items: center;
    color: white;

    .modal-title-button__wrapper {
      display: flex;
      align-items: center;
    }

    .v-btn--minimize {
      pointer-events: auto;
    }
  }

  > .v-card-title,
  > div:not(.v-card-text) {
    flex-shrink: 0;
  }

  .v-card-text {
    flex: 1 1 auto;
    min-height: 0;
    overflow-y: auto;
  }
}
</style>
