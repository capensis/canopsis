<template lang="pug">
  v-speed-dial.d-inline-block(
    v-model="isVSpeedDialOpen",
    direction="bottom",
    transition="scale-transition"
  )
    template(#activator="")
      v-tooltip(top)
        template(#activator="{ on }")
          v-btn.primary(v-on="on", :input-value="isVSpeedDialOpen", dark, fab, small)
            v-icon add
            v-icon close
        span {{ $t('context.fab.common') }}
    v-tooltip(bottom)
      template(#activator="{ on }")
        v-btn(
          v-on="on",
          color="indigo",
          fab,
          dark,
          small,
          @click.prevent.stop="showCreateServiceModal"
        )
          v-icon(size="24") $vuetify.icons.engineering
      span {{ $t('context.fab.addService') }}
</template>

<script>
import { MODALS } from '@/constants';

import entitiesServiceMixin from '@/mixins/entities/service';
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
