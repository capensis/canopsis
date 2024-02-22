<template>
  <c-speed-dial
    class="d-inline-block"
    direction="bottom"
    transition="scale-transition"
  >
    <template #activator="{ bind: speedDialBind }">
      <v-tooltip top>
        <template #activator="{ on: tooltipOn }">
          <v-btn
            v-bind="speedDialBind"
            class="primary"
            dark
            fab
            small
            v-on="tooltipOn"
          >
            <v-icon>add</v-icon>
            <v-icon>close</v-icon>
          </v-btn>
        </template>
        <span>{{ $t('context.fab.common') }}</span>
      </v-tooltip>
    </template>
    <v-tooltip bottom>
      <template #activator="{ on }">
        <v-btn
          color="indigo"
          fab
          dark
          small
          v-on="on"
          @click.prevent.stop="showCreateServiceModal"
        >
          <v-icon size="24">
            $vuetify.icons.engineering
          </v-icon>
        </v-btn>
      </template>
      <span>{{ $t('context.fab.addService') }}</span>
    </v-tooltip>
  </c-speed-dial>
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
