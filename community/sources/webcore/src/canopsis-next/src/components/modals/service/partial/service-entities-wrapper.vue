<template lang="pug">
  div
    v-layout.ma-2(v-if="selectedEntities.length", align-center, row)
      span.font-weight-bold {{ $t('serviceWeather.massActions') }}:
      service-entity-actions(:actions="actions", @apply="applyActionForSelected")
    div.mt-2(v-for="serviceEntity in serviceEntities", :key="serviceEntity.key")
      service-entity(
        :service-id="service._id",
        :entity="serviceEntity",
        :last-action-unavailable="hasActionError(serviceEntity)",
        :entity-name-field="entityNameField",
        :widget-parameters="widgetParameters",
        :selected="isEntitySelected(serviceEntity)",
        @select="updateSelected(serviceEntity, $event)",
        @remove-unavailable="removeEntityFromUnavailable(serviceEntity)",
        @add:action="$listeners['add:action']"
      )
</template>

<script>
import { getAvailableActionsByEntities, isActionTypeAvailableForEntity } from '@/helpers/entities/context';
import { filterById, mapIds } from '@/helpers/entities';

import { widgetActionPanelServiceEntityMixin } from '@/mixins/widget/actions-panel/service-entity';

import ServiceEntityActions from '@/components/modals/service/partial/service-entity-actions.vue';

import ServiceEntity from './service-entity.vue';

export default {
  inject: ['$actionsQueue'],
  components: {
    ServiceEntityActions,
    ServiceEntity,
  },
  mixins: [widgetActionPanelServiceEntityMixin],
  props: {
    service: {
      type: Object,
      required: true,
    },
    serviceEntities: {
      type: Array,
      default: () => [],
    },
    entityNameField: {
      type: String,
      default: 'name',
    },
    widgetParameters: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      selectedEntities: [],
    };
  },
  computed: {
    actions() {
      return getAvailableActionsByEntities(this.selectedEntities)
        .filter(this.actionsAccessFilterHandler)
        .map(action => ({
          ...action,
          disabled: this.isActionAvailableForSelectedAction(action.type),
        }));
    },

    selectedEntitiesIds() {
      return mapIds(this.selectedEntities);
    },

    pendingEntitiesIdsByActionType() {
      return this.$actionsQueue.queue
        .reduce((acc, { actionType, entities }) => {
          const entitiesIds = mapIds(entities);
          const pendingIds = acc[actionType];

          acc[actionType] = pendingIds
            ? pendingIds.concat(entitiesIds)
            : entitiesIds;

          return acc;
        }, {});
    },
  },
  methods: {
    hasActionError(entity) {
      return this.unavailableEntitiesAction[entity._id];
    },

    isEntityActionPending(type, id) {
      const pendingEntitiesIds = this.pendingEntitiesIdsByActionType[type] || [];

      return pendingEntitiesIds.includes(id);
    },

    isActionAvailableForSelectedAction(type) {
      return this.selectedEntities
        .filter(entity => isActionTypeAvailableForEntity(type, entity))
        .every(({ _id: id }) => this.isEntityActionPending(type, id));
    },

    applyActionForSelected({ type }) {
      const availableEntities = this.selectedEntities
        .filter(entity => !this.isEntityActionPending(type, entity._id));

      this.addEntityAction(type, availableEntities);
    },

    updateSelected(entity, checked) {
      if (checked) {
        this.selectedEntities.push(entity);
      } else {
        this.selectedEntities = filterById(this.selectedEntities, entity);
      }
    },

    isEntitySelected(entity) {
      return this.selectedEntitiesIds.includes(entity._id);
    },
  },
};
</script>
