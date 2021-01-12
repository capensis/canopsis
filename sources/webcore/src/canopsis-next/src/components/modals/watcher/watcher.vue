<template lang="pug">
  modal-wrapper(:title-color="color", close)
    template(slot="title")
      span {{ watcher.name }}
    template(slot="text")
      v-fade-transition(mode="out-in")
        watcher-template(
          v-if="!watcherEntitiesPending",
          :watcher="watcher",
          :watcher-entities="watcherEntitiesWithKey",
          :modal-template="config.modalTemplate",
          :entity-template="config.entityTemplate",
          :items-per-page="config.itemsPerPage",
          @add:event="addEventToQueue"
        )
        v-layout(v-else, column)
          v-flex(xs12)
            v-layout(justify-center)
              v-progress-circular(indeterminate, color="primary")
    template(slot="actions")
      v-alert.ma-0.pa-1.pr-2(
        :value="eventsQueue.length",
        type="info"
      ) {{ eventsQueue.length }} {{ $t('modals.watcher.actionPending') }}
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-tooltip.mx-2(top)
        v-btn.secondary(slot="activator", @click="fetchWatchersList")
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

import { MODALS, EVENT_ENTITY_TYPES, PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { formToPbehavior, pbehaviorToRequest } from '@/helpers/forms/planning-pbehavior';
import { addKeyInEntity } from '@/helpers/entities';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';
import eventActionsMixin from '@/mixins/event-actions/alarm';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';
import entitiesWatcherEntityMixin from '@/mixins/entities/watcher-entity';

import ModalWrapper from '../modal-wrapper.vue';
import ModalTitleButtons from '../modal-title-buttons.vue';

import WatcherTemplate from './partial/watcher-template.vue';

export default {
  name: MODALS.watcher,
  components: { ModalTitleButtons, WatcherTemplate, ModalWrapper },
  inject: ['$system'],
  mixins: [
    modalInnerMixin,
    eventActionsMixin,
    entitiesPbehaviorMixin,
    entitiesWatcherEntityMixin,
    submittableMixin(),
    confirmableModalMixin({ field: 'eventsQueue' }),
  ],
  data() {
    return {
      attributes: {},
      eventsQueue: [],
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
    this.fetchWatchersList();

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
    fetchWatchersList() {
      const params = {};

      if (this.modal.config.sort) {
        const { column, order } = this.modal.config.sort;

        params.sort_by = column;
        params.sort = order.toLowerCase();
      }

      this.fetchWatcherEntitiesList({
        watcherId: this.watcher._id,
        params,
      });
    },

    addEventToQueue(event) {
      this.eventsQueue.push(event);
    },

    getPausedPbehaviors(pbehaviors) {
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
      const requests = this.eventsQueue.reduce((acc, event) => {
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
  },
};
</script>
