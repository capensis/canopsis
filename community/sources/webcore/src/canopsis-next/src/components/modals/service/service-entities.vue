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
          @add:event="addEventToQueue",
          @refresh="fetchList"
        )
        v-layout(v-else, column)
          v-flex(xs12)
            v-layout(justify-center)
              v-progress-circular(color="primary", indeterminate)
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
import { MODALS, EVENT_ENTITY_TYPES, SORT_ORDERS } from '@/constants';
import { PAGINATION_LIMIT } from '@/config';

import { formToPbehavior, pbehaviorToRequest } from '@/helpers/forms/planning-pbehavior';
import { addKeyInEntities } from '@/helpers/entities';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';
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
    modalInnerMixin,
    eventActionsMixin,
    entitiesPbehaviorMixin,
    entitiesServiceEntityMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator({ field: 'events' }),
    localQueryMixin,
  ],
  data() {
    return {
      events: {
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

    addEventToQueue(event) {
      this.events.queue.push(event);
    },

    async submit() {
      const requests = this.events.queue.reduce((acc, event) => {
        if (event.type === EVENT_ENTITY_TYPES.pause) {
          const pbehavior = pbehaviorToRequest(formToPbehavior(event.data));

          acc.push(this.createPbehavior({ data: pbehavior }));
        } else if (event.type === EVENT_ENTITY_TYPES.play) {
          acc.push(this.removePbehavior({ id: event.data.pbehavior_info.id }));
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
