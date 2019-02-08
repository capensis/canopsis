<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ watcher.display_name }}
        v-btn(icon, dark, @click.native="hideModal")
          v-icon close
    v-divider
    v-card-text
      v-fade-transition
        div(v-show="!watcherEntitiesPending")
          watcher-template(
          :watcher="watcher",
          :watcherEntities="watcherEntities",
          :modalTemplate="config.modalTemplate",
          :entityTemplate="config.entityTemplate",
          @addEvent="addEventToQueue"
          )
      v-fade-transition
        v-layout(v-show="watcherEntitiesPending", column)
          v-flex(xs12)
            v-layout(justify-center)
              v-progress-circular(indeterminate, color="primary")
    v-divider
    v-layout.py-1(justify-end, align-center)
      v-alert.ma-0.pa-1.pr-2(
      :value="eventsQueue.length",
      type="info",
      ) {{ eventsQueue.length }} {{ $t('modals.watcher.actionPending') }}
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click="submit", :loading="submitting", :disabled="submitting") {{ $t('common.submit') }}
</template>

<script>
import { pick, mapValues } from 'lodash';

import VRuntimeTemplate from 'v-runtime-template';

import { MODALS, ENTITIES_TYPES, EVENT_ENTITY_TYPES, PBEHAVIOR_TYPES } from '@/constants';

import weatherEventMixin from '@/mixins/weather-event-actions';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';

import entitiesWatcherMixin from '@/mixins/entities/watcher';
import entitiesWatcherEntityMixin from '@/mixins/entities/watcher-entity';
import modalInnerMixin from '@/mixins/modal/inner';

import WatcherTemplate from './partial/watcher-template.vue';

export default {
  name: MODALS.watcher,
  components: {
    VRuntimeTemplate,
    WatcherTemplate,
  },
  mixins: [
    weatherEventMixin,
    entitiesPbehaviorMixin,
    modalInnerMixin,
    entitiesWatcherMixin,
    entitiesWatcherEntityMixin,
  ],
  data() {
    return {
      attributes: {},
      eventsQueue: [],
      submitting: false,
    };
  },
  computed: {
    watcher() {
      return this.getWatcher(this.config.watcherId);
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

    this.fetchWatcherEntitiesList({ watcherId: this.config.watcherId });
  },
  methods: {
    addEventToQueue(event) {
      this.eventsQueue.push(event);
    },

    async submit() {
      this.submitting = true;

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
              accSecond.push(this.removePbehavior({ id: pbehavior._id }));
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

      this.submitting = false;
      this.hideModal();
    },
  },
};
</script>
