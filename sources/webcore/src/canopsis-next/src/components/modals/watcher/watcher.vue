<template lang="pug">
  modal-wrapper
    template(slot="fullTitle")
      v-card-title.white--text(:style="{ backgroundColor: color }")
        v-layout(justify-space-between, align-center)
          span.headline {{ watcher.display_name }}
    template(slot="text")
      v-fade-transition
        div(v-show="!watcherEntitiesPending")
          watcher-template(
            :watcher="watcher",
            :watcherEntities="watchers",
            :watchersMeta="metaData",
            :modalTemplate="config.modalTemplate",
            :entityTemplate="config.entityTemplate",
            @addEvent="addEventToQueue",
            @change:page="changePage",
            @change:limit="changeLimit"
          )
      v-fade-transition
        v-layout(v-show="watcherEntitiesPending", column)
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
        v-btn.secondary(slot="activator", @click="refresh")
          v-icon refresh
        span {{ $t('modals.watcher.refreshEntities') }}
      v-btn.primary(
        :disabled="isDisabled",
        :loading="submitting",
        @click="submit"
      ) {{ $t('common.submit') }}
</template>

<script>
import { pick, mapValues } from 'lodash';

import { MODALS, ENTITIES_TYPES, EVENT_ENTITY_TYPES, PBEHAVIOR_TYPES } from '@/constants';

import watcherQueryMixin from '@/mixins/watcher/query';
import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import eventActionsMixin from '@/mixins/event-actions/alarm';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';

import ModalWrapper from '../modal-wrapper.vue';

import WatcherTemplate from './partial/watcher-template.vue';

export default {
  name: MODALS.watcher,
  components: { WatcherTemplate, ModalWrapper },
  mixins: [
    modalInnerMixin,
    watcherQueryMixin,
    submittableMixin(),
    eventActionsMixin,
    entitiesPbehaviorMixin,
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
  },
  mounted() {
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
    addEventToQueue(event) {
      this.eventsQueue.push(event);
    },

    refresh() {
      this.fetchWatcherEntitiesList({ watcherId: this.watcher.entity_id });
    },

    async submit() {
      const requests = this.eventsQueue.reduce((acc, event) => {
        if (event.type === EVENT_ENTITY_TYPES.pause) {
          acc.push(this.createPbehavior({
            data: event.data,
            parents: [event.entity],
            parentsType: ENTITIES_TYPES.entity,
          }));
        } else if (event.type === EVENT_ENTITY_TYPES.play) {
          const pausedPbehaviorsRequests = event.data.pbehavior.reduce((accSecond, pbehavior) => {
            if (pbehavior.type_ === PBEHAVIOR_TYPES.pause) {
              const data = {
                ...pick(pbehavior, ['author', 'exdate', 'filter', 'name', 'reason', 'rrule', 'tstart', 'type_']),
                tstop: Math.round(Date.now() / 1000),
              };

              accSecond.push(this.updatePbehavior({ data, id: pbehavior._id }));
            }

            return accSecond;
          }, []);

          acc.push(...pausedPbehaviorsRequests);
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
