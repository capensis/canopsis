<template>
  <v-speed-dial
    class="d-inline-block"
    v-model="isVSpeedDialOpen"
    direction="bottom"
    transition="scale-transition"
  >
    <template #activator="">
      <v-tooltip top="top">
        <template #activator="{ on }">
          <v-btn
            class="primary"
            v-on="on"
            :input-value="isVSpeedDialOpen"
            dark="dark"
            fab="fab"
            small="small"
          >
            <v-icon>add</v-icon>
            <v-icon>close</v-icon>
          </v-btn>
        </template><span>{{ $t('context.fab.common') }}</span>
      </v-tooltip>
    </template>
    <v-tooltip bottom="bottom">
      <template #activator="{ on }">
        <v-btn
          v-on="on"
          color="indigo"
          fab="fab"
          dark="dark"
          small="small"
          @click.prevent.stop="showCreateServiceModal"
        >
          <v-icon size="24">
            $vuetify.icons.engineering
          </v-icon>
        </v-btn>
      </template><span>{{ $t('context.fab.addService') }}</span>
    </v-tooltip>
  </v-speed-dial>
</template>

<script>
import { MODALS } from '@/constants';

import { entitiesServiceMixin } from '@/mixins/entities/service';
import { entitiesContextEntityMixin } from '@/mixins/entities/context-entity';

/**
 * Buttons to open the modal to add entities
 *
 * @module context
 */
export default {
  mixins: [
    entitiesServiceMixin,
    entitiesContextEntityMixin,
  ],
  data() {
    return {
      isVSpeedDialOpen: false,
    };
  },
  methods: {
    showCreateServiceModal() {
      this.$modals.show({
        name: MODALS.createService,
        config: {
          title: this.$t('modals.createService.create.title'),
          action: async (service) => {
            await this.createService({ data: service });

            this.$popups.success({ text: this.$t('modals.createService.success.create') });
          },
        },
      });
    },
  },
};
</script>
