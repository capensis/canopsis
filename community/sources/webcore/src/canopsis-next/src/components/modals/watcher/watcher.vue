<template lang="pug">
  modal-wrapper(:title-color="color", close)
    template(slot="title")
      span {{ watcher.name }}
    template(slot="text")
      v-fade-transition(mode="out-in")
        watcher-template(
          v-if="!watcherEntitiesPendingOnMount",
          :watcher="watcher",
          :watcher-entities="watcherEntitiesWithKey",
          :modal-template="config.modalTemplate",
          :entity-template="config.entityTemplate",
          :pagination.sync="pagination",
          :total-items="watcherEntitiesMeta.total_count",
          :pending="watcherEntitiesPending",
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
          span {{ events.queue.length }} {{ $t('modals.watcher.actionPending') }}
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-tooltip.mx-2(top)
        v-btn.secondary(slot="activator", @click="fetchList")
          v-icon refresh
        span {{ $t('modals.watcher.refreshEntities') }}
      v-btn.primary(
        :disabled="isDisabled",
        :loading="submitting",
        @click="submit"
      ) {{ $t('common.submit') }}
</template>

<script>
import moment from 'moment-timezone';
import { pick, mapValues } from 'lodash';

import { MODALS, EVENT_ENTITY_TYPES, PBEHAVIOR_TYPE_TYPES, SORT_ORDERS } from '@/constants';

import { formToPbehavior, pbehaviorToRequest } from '@/helpers/forms/planning-pbehavior';
import { addKeyInEntity } from '@/helpers/entities';

import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';
import eventActionsMixin from '@/mixins/event-actions/alarm';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';
import entitiesWatcherEntityMixin from '@/mixins/entities/watcher-entity';
import localQueryMixin from '@/mixins/query-local/query';

import ModalWrapper from '../modal-wrapper.vue';
import ModalTitleButtons from '../modal-title-buttons.vue';

import WatcherTemplate from './partial/watcher-template.vue';

export default {
  name: MODALS.watcher,
  inject: ['$system'],
  provide() {
    return {
      $eventsQueue: this.events,
    };
  },
  components: { ModalTitleButtons, WatcherTemplate, ModalWrapper },
  mixins: [
    eventActionsMixin,
    entitiesPbehaviorMixin,
    entitiesWatcherEntityMixin,
    submittableMixin(),
    confirmableModalMixin({ field: 'events' }),
    localQueryMixin,
  ],
  data() {
    return {
      attributes: {},
      events: {
        queue: [],
      },
      watcherEntitiesPendingOnMount: false,
      query: {
        rowsPerPage: this.modal.config.itemsPerPage,
        sortKey: 'state',
        sortDir: SORT_ORDERS.desc,
      },
    };
  },
  computed: {
    watcher() {
      return this.config.watcher;
    },

    color() {
      return this.config.color;
    },

    watcherEntitiesWithKey() {
      return addKeyInEntity(this.watcherEntities);
    },
  },
  mounted() {
    this.watcherEntitiesPendingOnMount = true;

    this.fetchList();

    this.watcherEntitiesPendingOnMount = false;

    // TODO: Do we need it ?
    const infoAttributes = mapValues(pick(this.watcher.infos, [
      'application_crit_label',
      'product_line',
      'service_period',
      'isInCarat',
      'application_label',
      'target_platform',
    ]), v => v.value);

    this.attributes = {
      org: this.watcher.org,
      ...infoAttributes,
    };
  },
  methods: {
    fetchList() {
      this.fetchWatcherEntitiesList({
        watcherId: this.watcher._id,
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
