<template>
  <v-form @submit.prevent>
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.createMap.title') }}</span>
      </template>

      <template #text>
        <c-layout column>
          <v-col
            v-for="type in availableTypes"
            :key="type"
            class="my-1 pa-0 cursor-pointer"
            @click="selectType(type)"
          >
            <v-card>
              <v-card-title class="py-3">
                <c-layout align-center>
                  <div class="text-subtitle-1">
                    {{ t(`map.types.${type}`) }}
                  </div>
                  <v-spacer />
                  <v-icon color="secondary">
                    {{ getIconByType(type) }}
                  </v-icon>
                </c-layout>
              </v-card-title>
            </v-card>
          </v-col>
        </c-layout>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { MODALS, MAP_TYPES, MAP_ICON_BY_TYPES, CREATE_MAP_MODAL_NAMES_BY_TYPE } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

export default {
  name: MODALS.createMap,

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();

    const availableTypes = computed(() => Object.values(MAP_TYPES));

    const getIconByType = (type) => MAP_ICON_BY_TYPES[type];

    const selectType = (type) => {
      modals.show({
        name: CREATE_MAP_MODAL_NAMES_BY_TYPE[type],
        config: {
          map: { type },
          action: async (map) => {
            await props.modal.config.action?.(map);

            modals.hide({ id: props.modal.id });
          },
        },
      });
    };

    return {
      t,
      availableTypes,
      getIconByType,
      selectType,
    };
  },
};
</script>
