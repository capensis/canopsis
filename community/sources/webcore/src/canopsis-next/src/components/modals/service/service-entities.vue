<template lang="pug">
  form(@submit.prevent="submit")
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
                :pending="serviceEntitiesPending",
                @refresh="fetchList",
                @add:action="addActionToQueue"
              )
              v-layout.pa-4(v-else, justify-center)
                v-progress-circular(color="primary", indeterminate)
          v-tab(:disabled="!hasPbehaviorListAccess") {{ $tc('common.activePbehavior') }}
          v-tab-item(lazy)
            pbehaviors-simple-list(
              :entity="service",
              with-active-status,
              creatable,
              editable,
              deletable,
              dense
            )
      template(#actions="")
        v-alert.ma-0.pa-1.pr-2(:value="!!actionsCount", color="info")
          v-layout(row, align-center)
            v-btn.mr-2(icon, small, @click="clearActions")
              v-icon(color="white", small) close
            span {{ actionsCount }} {{ $tc('modals.service.actionPending', actionsCount) }}
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-tooltip.mx-2(top)
          template(#activator="{ on }")
            v-btn.secondary(v-on="on", @click="fetchList")
              v-icon refresh
          span {{ $t('modals.service.refreshEntities') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { pick, isEmpty } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';

import { MODALS, SORT_ORDERS, USERS_PERMISSIONS, WEATHER_ACTIONS_TYPES } from '@/constants';

import { addKeyInEntities } from '@/helpers/entities';
import { createDowntimePbehavior } from '@/helpers/entities/pbehavior';
import { convertActionsToEvents } from '@/helpers/entities/entity';

import { authMixin } from '@/mixins/auth';
import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';
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
  provide() {
    return {
      $actionsQueue: this.actions,
    };
  },
  components: { PbehaviorsSimpleList, ServiceTemplate, ModalWrapper },
  mixins: [
    authMixin,
    localQueryMixin,
    modalInnerMixin,
    eventActionsAlarmMixin,
    eventActionsMixin,
    entitiesPbehaviorMixin,
    entitiesServiceEntityMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator({ field: 'actions' }),
  ],
  data() {
    return {
      actions: {
        queue: [],
      },
      pending: true,
      query: {
        rowsPerPage: this.modal.config.widgetParameters.modalItemsPerPage || PAGINATION_LIMIT,
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

    actionsCount() {
      return this.actions.queue.reduce((count, { entities }) => count + entities.length, 0);
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

    addActionToQueue(action) {
      this.actions.queue.push(action);
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

    async submit() {
      const { eventsActions, createdPbehaviors, removedPbehaviors } = this.actions.queue.reduce((acc, action) => {
        if (action.actionType === WEATHER_ACTIONS_TYPES.entityPause) {
          acc.createdPbehaviors.push(
            ...this.getCreatedPbehaviorsByEntitites(action.entities, action.payload),
          );
        } else if (action.actionType === WEATHER_ACTIONS_TYPES.entityPlay) {
          acc.removedPbehaviors.push(
            ...this.getPausedPbehaviorsByEntitites(action.entities),
          );
        } else {
          acc.eventsActions.push(action);
        }

        return acc;
      }, {
        createdPbehaviors: [],
        removedPbehaviors: [],
        eventsActions: [],
      });

      const events = convertActionsToEvents(eventsActions);

      await Promise.all([
        createdPbehaviors.length && this.createPbehaviorsWithPopups(createdPbehaviors),
        removedPbehaviors.length && this.removePbehaviorsWithPopups(removedPbehaviors),
        events.length && this.createEventAction({ data: events }),
      ]);

      this.$modals.hide();
    },

    clearActions() {
      this.actions.queue = [];
    },
  },
};
</script>
