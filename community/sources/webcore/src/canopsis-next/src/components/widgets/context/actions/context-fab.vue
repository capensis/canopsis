<template lang="pug">
  v-speed-dial.d-inline-block(
    v-model="isVSpeedDialOpen",
    direction="left",
    transition="scale-transition"
  )
    v-btn.primary(slot="activator", :input-value="isVSpeedDialOpen", dark, fab, small)
      v-icon add
      v-icon close
    v-tooltip(top)
      v-btn(
        slot="activator",
        color="indigo",
        fab,
        dark,
        small,
        @click.prevent.stop="showCreateServiceModal"
      )
        v-icon(size="24") $vuetify.icons.engineering
      span {{ $t('common.service') }}
</template>

<script>
import { MODALS } from '@/constants';

import entitiesServiceMixin from '@/mixins/entities/service';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';

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
