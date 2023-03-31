<template lang="pug">
  modal-wrapper(:title-color="color", close)
    template(#title="")
      span {{ service.name }}
    template(#text="")
      v-tabs.position-relative(slider-color="primary", fixed-tabs)
        v-tab {{ $t('common.service') }}
        v-tab-item
          c-progress-overlay(:pending="pending")
          service-template(
            :service="service",
            :service-entities="serviceEntitiesWithKey",
            :widget-parameters="widgetParameters",
            :pagination.sync="pagination",
            :total-items="serviceEntitiesMeta.total_count",
            @refresh="refresh",
            @apply:action="applyAction"
          )
        v-tab(:disabled="!hasPbehaviorListAccess") {{ $tc('common.pbehavior', 2) }}
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
          v-btn.secondary(v-on="on", @click="refresh")
            v-icon refresh
        span {{ $t('modals.service.refreshEntities') }}
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.close') }}
</template>

<script>
import { pick, isEmpty } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';

import { MODALS, PBEHAVIOR_ORIGINS, SORT_ORDERS, USERS_PERMISSIONS, WEATHER_ACTIONS_TYPES } from '@/constants';

import Observer from '@/services/observer';

import { createDowntimePbehavior } from '@/helpers/entities/pbehavior';
import { convertActionToEvents } from '@/helpers/entities/entity';

import { authMixin } from '@/mixins/auth';
import { modalInnerMixin } from '@/mixins/modal/inner';
import { eventActionsAlarmMixin } from '@/mixins/event-actions/alarm';
import { eventActionsMixin } from '@/mixins/event-actions';
import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';
import { entitiesServiceEntityMixin } from '@/mixins/entities/service-entity';
import { localQueryMixin } from '@/mixins/query-local/query';

import ServiceTemplate from '@/components/other/service/partials/service-template.vue';
import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/partials/pbehaviors-simple-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.serviceEntities,
  provide() {
    return {
      $periodicRefresh: this.$periodicRefresh,
    };
  },
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
      return this.serviceEntities.map(entity => ({
        ...entity,
        key: `${entity._id}_${entity.alarm_id}`,
      }));
    },

    hasPbehaviorListAccess() {
      return this.checkAccess(USERS_PERMISSIONS.business.serviceWeather.actions.pbehaviorList);
    },
  },

  beforeCreate() {
    this.$periodicRefresh = new Observer();
  },

  created() {
    this.$periodicRefresh.register(this.fetchList);
  },

  beforeDestroy() {
    this.$periodicRefresh.unregister(this.fetchList);
  },

  async mounted() {
    await this.fetchList();
  },

  methods: {
    refresh() {
      return this.$periodicRefresh.notify();
    },

    async fetchList() {
      this.pending = true;

      const params = this.getQuery();
      params.with_instructions = true;
      params.pbh_origin = PBEHAVIOR_ORIGINS.serviceWeather;

      await this.fetchServiceEntitiesList({ id: this.service._id, params });

      this.pending = false;
    },

    async createPauseEvent(entities, payload) {
      const response = await this.createEntityPbehaviors({
        data: entities.map(entity => createDowntimePbehavior({
          entity,
          ...pick(payload, ['comment', 'reason', 'type']),
        }), []),
      });

      this.showPbehaviorResponseErrorPopups(response);
    },

    async createPlayEvent(entities) {
      const response = await this.removeEntityPbehaviors({
        data: entities.map(({ _id: id }) => ({
          entity: id,
          origin: PBEHAVIOR_ORIGINS.serviceWeather,
        })),
      });

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

    async applyAction({ entities, actionType, payload }) {
      if (actionType === WEATHER_ACTIONS_TYPES.entityPause) {
        await this.createPauseEvent(entities, payload);
      } else if (actionType === WEATHER_ACTIONS_TYPES.entityPlay) {
        await this.createPlayEvent(entities);
      } else {
        const events = entities.reduce((acc, entity) => {
          acc.push(...convertActionToEvents({
            entity,
            actionType,
            payload,
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
