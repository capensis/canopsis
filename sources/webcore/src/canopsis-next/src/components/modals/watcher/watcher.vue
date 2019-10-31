<template lang="pug">
  v-card
    v-card-title.white--text(:style="{ backgroundColor: color }")
      v-layout(justify-space-between, align-center)
        span.headline {{ watcher.display_name }}
        v-btn(icon, dark, @click.native="$modals.hide")
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
        type="info"
      ) {{ eventsQueue.length }} {{ $t('modals.watcher.actionPending') }}
      v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
      v-tooltip(top)
        v-btn(@click="refresh", color="secondary", slot="activator")
          v-icon refresh
        span {{ $t('modals.watcher.refreshEntities') }}
      v-btn.primary(@click="submit", :loading="submitting", :disabled="submitting") {{ $t('common.submit') }}
</template>

<script>
import { pick, mapValues } from 'lodash';

import { MODALS, ENTITIES_TYPES, EVENT_ENTITY_TYPES, PBEHAVIOR_TYPES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import eventActionsMixin from '@/mixins/event-actions/alarm';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';
import entitiesWatcherEntityMixin from '@/mixins/entities/watcher-entity';

import WatcherTemplate from './partial/watcher-template.vue';

export default {
  name: MODALS.watcher,
  components: { WatcherTemplate },
  mixins: [
    modalInnerMixin,
    eventActionsMixin,
    entitiesPbehaviorMixin,
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

    this.fetchWatcherEntitiesList({ watcherId: this.watcher.entity_id });
  },
  methods: {
    addEventToQueue(event) {
      this.eventsQueue.push(event);
    },

    refresh() {
      this.fetchWatcherEntitiesList({ watcherId: this.watcher.entity_id });
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
      this.$modals.hide();
    },
  },
};
</script>
