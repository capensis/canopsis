<template lang="pug">
  modal-wrapper(:title-color="color", close)
    template(#title="")
      span {{ service.name }}
    template(#text="")
      v-tabs(slider-color="primary", fixed-tabs, light)
        v-tab {{ $t('common.service') }}
        v-tab-item
          v-fade-transition(mode="out-in")
            service-template(
              v-if="!pending",
              :service="service",
              :service-entities="serviceEntitiesWithKey",
              :widget-parameters="widgetParameters",
              :pagination.sync="pagination",
              :total-items="serviceEntitiesMeta.total_count",
              :unavailable-entities-action="unavailableEntitiesAction",
              :pending="serviceEntitiesPending",
              @refresh="fetchList",
              @remove:unavailable="removeEntityFromUnavailable",
              @apply:action="applyAction"
            )
            v-layout.pa-4(v-else, justify-center)
              v-progress-circular(color="primary", indeterminate)
        v-tab(:disabled="!hasPbehaviorListAccess") {{ $tc('common.activePbehavior') }}
        v-tab-item(lazy)
          pbehaviors-simple-list(
            :entity="service",
            with-active-status,
            addable,
            updatable,
            removable,
            dense
          )
    template(#actions="")
      v-tooltip.mx-2(top)
        template(#activator="{ on }")
          v-btn.secondary(v-on="on", @click="fetchList")
            v-icon refresh
        span {{ $t('modals.service.refreshEntities') }}
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.close') }}
</template>

<script>
import { pick, isEmpty } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';

import { MODALS, SORT_ORDERS, USERS_PERMISSIONS, WEATHER_ACTIONS_TYPES } from '@/constants';

import { addKeyInEntities } from '@/helpers/entities';
import { createDowntimePbehavior } from '@/helpers/entities/pbehavior';
import { convertActionToEvents, isActionTypeAvailableForEntity } from '@/helpers/entities/entity';

import { authMixin } from '@/mixins/auth';
import { modalInnerMixin } from '@/mixins/modal/inner';
import { eventActionsAlarmMixin } from '@/mixins/event-actions/alarm';
import { eventActionsMixin } from '@/mixins/event-actions';
import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';
import { entitiesServiceEntityMixin } from '@/mixins/entities/service-entity';
import { localQueryMixin } from '@/mixins/query-local/query';

import ServiceTemplate from '@/components/other/service/partials/service-template.vue';
import PbehaviorsSimpleList from '@/components/other/pbehavior/partials/pbehaviors-simple-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.serviceEntities,
  inject: ['$system'],
  components: { PbehaviorsSimpleList, ServiceTemplate, ModalWrapper },
  mixins: [
    authMixin,
    localQueryMixin,
    modalInnerMixin,
    eventActionsAlarmMixin,
    eventActionsMixin,
    entitiesPbehaviorMixin,
    entitiesServiceEntityMixin,
  ],
  data() {
    return {
      pending: true,
      unavailableEntitiesAction: {},
      query: {
        rowsPerPage: this.modal.config.widgetParameters.modalItemsPerPage ?? PAGINATION_LIMIT,
        sortKey: 'state',
        sortDir: SORT_ORDERS.desc,
      },
    };
  },
  computed: {
    service() {
      return this.config.service;
    },

    color() {
      return this.config.color;
    },

    widgetParameters() {
      return this.config.widgetParameters;
    },

    serviceEntitiesWithKey() {
      return addKeyInEntities(this.serviceEntities);
    },

    hasPbehaviorListAccess() {
      return this.checkAccess(USERS_PERMISSIONS.business.serviceWeather.actions.pbehaviorList);
    },
  },
  async mounted() {
    this.pending = true;

    await this.fetchList();

    this.pending = false;
  },
  methods: {
    fetchList() {
      const params = this.getQuery();

      params.with_instructions = true;

      return this.fetchServiceEntitiesList({
        id: this.service._id,
        params,
      });
    },

    getCreatedPbehaviorsByEntitites(entities, data) {
      return entities.reduce((acc, entity) => {
        acc.push(createDowntimePbehavior({
          entity,
          ...pick(data, ['comment', 'reason', 'type']),
        }));

        return acc;
      }, []);
    },

    getPausedPbehaviorsByEntitites(entities) {
      return entities.map(entity => ({ _id: entity.pbehavior_info.id }));
    },

    async createPbehaviorsWithPopups(pbehaviors) {
      const response = await this.createPbehaviorsWithComments(pbehaviors);

      this.showPbehaviorResponseErrorPopups(response);
    },

    async removePbehaviorsWithPopups(pbehaviors) {
      const response = await this.removePbehaviors(pbehaviors);

      this.showPbehaviorResponseErrorPopups(response);
    },

    showPbehaviorResponseErrorPopups(response) {
      if (response?.length) {
        response.forEach(({ error, errors }) => {
          if (error || !isEmpty(errors)) {
            this.$popups.error({ text: error || Object.values(errors).join('\n') });
          }
        });
      }
    },

    getAvailableEntities(action) {
      const {
        availableEntities,
        unavailableEntities,
      } = action.entities.reduce((acc, entity) => {
        if (isActionTypeAvailableForEntity(action.actionType, entity)) {
          acc.availableEntities.push(entity);
        } else {
          acc.unavailableEntities.push(entity);
        }

        return acc;
      }, {
        availableEntities: [],
        unavailableEntities: [],
      });

      this.unavailableEntitiesAction = unavailableEntities.reduce((acc, { _id: id }) => {
        acc[id] = true;

        return acc;
      }, {});

      return availableEntities;
    },

    removeEntityFromUnavailable(entity) {
      this.unavailableEntitiesAction[entity._id] = false;
    },

    async applyAction(action) {
      const availableEntities = this.getAvailableEntities(action);

      if (action.actionType === WEATHER_ACTIONS_TYPES.entityPause) {
        await this.createPbehaviorsWithPopups(
          this.getCreatedPbehaviorsByEntitites(availableEntities, action.payload),
        );
      } else if (action.actionType === WEATHER_ACTIONS_TYPES.entityPlay) {
        await this.removePbehaviorsWithPopups(
          this.getPausedPbehaviorsByEntitites(availableEntities),
        );
      } else {
        const events = availableEntities.reduce((acc, entity) => {
          acc.push(...convertActionToEvents({
            entity,
            actionType: action.actionType,
            payload: action.payload,
          }));

          return acc;
        }, []);

        await this.createEventAction({ data: events });
      }

      await this.fetchList();
    },
  },
};
</script>
