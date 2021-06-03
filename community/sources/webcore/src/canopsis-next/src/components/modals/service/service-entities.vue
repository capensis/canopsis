<template lang="pug">
  modal-wrapper(:title-color="color", close)
    template(slot="title")
      span {{ service.name }}
    template(slot="text")
      v-fade-transition(mode="out-in")
        service-template(
          v-if="!pending",
          :service="service",
          :service-entities="serviceEntitiesWithKey",
          :widget-parameters="widgetParameters",
          :pagination.sync="pagination",
          :total-items="serviceEntitiesMeta.total_count",
          :pending="serviceEntitiesPending",
          @add:event="addEventToQueue"
        )
        v-layout(v-else, column)
          v-flex(xs12)
            v-layout(justify-center)
              v-progress-circular(indeterminate, color="primary")
    template(slot="actions")
      v-alert.ma-0.pa-1.pr-2(
        :value="events.queue.length",
        color="info"
      )
        v-layout(row, align-center)
          v-btn.mr-2(icon, small, @click="clearActions")
            v-icon(color="white", small) close
          span {{ events.queue.length }} {{ $t('modals.service.actionPending') }}
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-tooltip.mx-2(top)
        v-btn.secondary(slot="activator", @click="fetchList")
          v-icon refresh
        span {{ $t('modals.service.refreshEntities') }}
      v-btn.primary(
        :disabled="isDisabled",
        :loading="submitting",
        @click="submit"
      ) {{ $t('common.submit') }}
</template>

<script>
import moment from 'moment-timezone';

import { MODALS, EVENT_ENTITY_TYPES, PBEHAVIOR_TYPE_TYPES, SORT_ORDERS } from '@/constants';

import { formToPbehavior, pbehaviorToRequest } from '@/helpers/forms/planning-pbehavior';
import { addKeyInEntities } from '@/helpers/entities';

import { submittableMixin } from '@/mixins/submittable';
import { confirmableModalMixin } from '@/mixins/confirmable-modal';
import eventActionsMixin from '@/mixins/event-actions/alarm';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';
import entitiesServiceEntityMixin from '@/mixins/entities/service-entity';
import { localQueryMixin } from '@/mixins/query-local/query';

import ModalWrapper from '../modal-wrapper.vue';

import ServiceTemplate from './partial/service-template.vue';

export default {
  name: MODALS.serviceEntities,
  inject: ['$system'],
  provide() {
    return {
      $eventsQueue: this.events,
    };
  },
  components: { ServiceTemplate, ModalWrapper },
  mixins: [
    eventActionsMixin,
    entitiesPbehaviorMixin,
    entitiesServiceEntityMixin,
    submittableMixin(),
    confirmableModalMixin({ field: 'events' }),
    localQueryMixin,
  ],
  data() {
    return {
      events: {
        queue: [],
      },
      pending: true,
      query: {
        rowsPerPage: this.modal.config.itemsPerPage,
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
  },
  async mounted() {
    this.pending = true;

    await this.fetchList();

    this.pending = false;
  },
  methods: {
    fetchList() {
      return this.fetchServiceEntitiesList({
        id: this.service._id,
        params: this.getQuery(),
      });
    },

    addEventToQueue(event) {
      this.events.queue.push(event);
    },

    getPausedPbehaviors(pbehaviors = []) {
      return pbehaviors.reduce((accSecond, pbehavior) => {
        if (pbehavior.type.type === PBEHAVIOR_TYPE_TYPES.pause) {
          accSecond.push(this.updatePbehavior({
            id: pbehavior._id,
            data: pbehaviorToRequest({
              ...pbehavior,

              tstop: moment().unix(),
            }),
          }));
        }

        return accSecond;
      }, []);
    },

    async submit() {
      const requests = this.events.queue.reduce((acc, event) => {
        if (event.type === EVENT_ENTITY_TYPES.pause) {
          const pbehavior = pbehaviorToRequest(formToPbehavior(event.data));

          acc.push(this.createPbehavior({ data: pbehavior }));
        } else if (event.type === EVENT_ENTITY_TYPES.play) {
          acc.push(...this.getPausedPbehaviors(event.data.pbehaviors));
        } else {
          acc.push(this.createEventAction({ data: event.data }));
        }

        return acc;
      }, []);

      await Promise.all(requests);

      this.$modals.hide();
    },

    clearActions() {
      this.events.queue = [];
    },
  },
};
</script>