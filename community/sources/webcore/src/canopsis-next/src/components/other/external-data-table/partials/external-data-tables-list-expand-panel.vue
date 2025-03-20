<template>
  <div class="secondary pa-3">
    <v-card>
      <v-card-text>
        <v-card>
          <v-card-text>
            <v-layout class="gap-4" column>
              <v-flex align-self-center>
                <span class="grey--text">{{ $t('common.noDataAvailable') }}</span>
              </v-flex>
              <v-flex align-self-center>
                <v-btn
                  color="primary"
                  @click="showImportExternalDataTablesModal"
                >
                  {{ $t('common.importData') }}
                </v-btn>
              </v-flex>
            </v-layout>
          </v-card-text>
        </v-card>
      </v-card-text>
    </v-card>
  </div>
</template>
<script>
import { MODALS } from '@/constants';

import { useModals } from '@/hooks/modals';

export default {
  props: {
    externalDataTable: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const modals = useModals();

    const showImportExternalDataTablesModal = () => modals.show({
      name: MODALS.importExternalDataTableRecords,
      config: {
        externalDataTable: props.externalDataTable,
      },
    });

    return {
      showImportExternalDataTablesModal,
    };
  },
};
</script>
