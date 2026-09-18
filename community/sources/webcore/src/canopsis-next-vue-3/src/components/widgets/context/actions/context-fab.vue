<template>
  <c-speed-dial
    class="d-inline-block"
    location="bottom center"
    transition="scale-transition"
  >
    <template #activator="{ props: activatorProps, isActive }">
      <v-tooltip location="top">
        <template #activator="{ props: tooltipProps }">
          <v-btn
            v-bind="{ ...tooltipProps, ...activatorProps }"
            :class="{ 'c-speed-dial__btn--active': isActive }"
            class="c-speed-dial__btn context-fab__btn"
            color="primary"
            size="x-small"
          >
            <v-icon icon="add" />
            <v-icon icon="close" />
          </v-btn>
        </template>
        <span>{{ $t('context.fab.common') }}</span>
      </v-tooltip>
    </template>
    <c-action-fab-btn
      :tooltip="$t('context.fab.addService')"
      color="indigo"
      icon="engineering"
      @click="showCreateServiceModal"
    />
  </c-speed-dial>
</template>

<script>
import { MODALS } from '@/constants';

import { entitiesContextEntityMixin } from '@/mixins/entities/context-entity';
import { entitiesServiceMixin } from '@/mixins/entities/service';

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

  methods: {
    showCreateServiceModal() {
      this.$modals.show({
        name: MODALS.createService,
        config: {
          title: this.$t('modals.createService.create.title'),
          action: async (service) => {
            const result = await this.createService({ data: service });

            this.$popups.success({ text: this.$t('modals.createService.success.create') });

            await this.fetchContextEntitiesListWithPreviousParams();

            return result;
          },
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.context-fab__btn {
  --v-btn-height: 40px;

  width: 40px;
  height: 40px;
  border-radius: 50%;
}
</style>
