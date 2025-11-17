<template>
  <v-chip
    :small="small"
    :color="color"
    :outlined="outlined"
    class="pl-2 pr-1"
  >
    <slot />
    <c-simple-tooltip :content="tooltipText" top>
      <template #activator="{ on }">
        <v-btn
          icon
          small
          @click.prevent.stop="showRelationsModal"
          v-on="on"
        >
          <v-icon
            :color="iconColor"
            size="16"
          >
            $vuetify.icons.flow
          </v-icon>
        </v-btn>
      </template>
    </c-simple-tooltip>
  </v-chip>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

export default {
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    small: {
      type: Boolean,
      default: false,
    },
    color: {
      type: String,
      required: false,
    },
    iconColor: {
      type: String,
      default: 'white',
    },
    outlined: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const entity = computed(() => props.alarm.entity);
    const hasUpstream = computed(() => entity.value.upstream !== undefined);

    const title = computed(() => (hasUpstream.value
      ? t('modals.entityUpstream.topLevelEntities')
      : t('modals.entityUpstream.entities')));

    const tooltipText = computed(() => (hasUpstream.value
      ? t('modals.entityUpstream.seeTopEntities')
      : t('modals.entityUpstream.seeEntities')));

    /**
     * Show entity relations modal with upstream and downstream network graph
     */
    const showRelationsModal = () => {
      modals.show({
        name: MODALS.entityUpstream,
        config: {
          entity: { ...entity.value, status: props.alarm.v.status.val },
          title: title.value,
        },
      });
    };

    return {
      tooltipText,
      showRelationsModal,
    };
  },
};
</script>
