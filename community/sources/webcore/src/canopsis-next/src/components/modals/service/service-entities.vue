<template>
  <modal-wrapper
    :title-color="color"
    close
  >
    <template #title="">
      <span>{{ service.name }}</span>
    </template>
    <template #text="">
      <v-tabs
        class="position-relative"
        slider-color="primary"
        centered
      >
        <v-tab>{{ $t('common.service') }}</v-tab>
        <v-tab-item>
          <c-progress-overlay :pending="pending" />
          <service-template
            :service="service"
            :service-entities="serviceEntitiesWithKey"
            :widget-parameters="widgetParameters"
            :options.sync="options"
            :total-items="serviceEntitiesMeta.total_count"
            :actions-requests="actionsRequests"
            @refresh="refresh"
            @add:action="addAction"
          />
        </v-tab-item>
        <v-tab :disabled="!hasPbehaviorListAccess">
          {{ $tc('common.pbehavior', 2) }}
        </v-tab>
        <v-tab-item>
          <pbehaviors-simple-list
            :entity="service"
            with-active-status
            addable
            updatable
            removable
            dense
          />
        </v-tab-item>
        <template v-if="hasAccessToCommentsList">
          <v-tab>{{ $tc('common.comment', 2) }}</v-tab>
          <v-tab-item>
            <entity-comments
              :entity="service"
              :addable="hasAccessToCreateComment"
              :editable="hasAccessToEditComment"
            />
          </v-tab-item>
        </template>
      </v-tabs>
    </template>
    <template #actions="">
      <v-alert
        :value="actionsRequests.length > 0"
        class="actions-requests-alert my-0 mr-2 pa-1 pr-2"
        type="info"
        dismissible
        @input="clearActions"
      >
        {{ actionsRequests.length }} {{ $tc('modals.service.actionInQueue', actionsRequests.length) }}
      </v-alert>
      <v-tooltip top>
        <template #activator="{ on }">
          <v-btn
            class="mx-2"
            color="secondary"
            v-on="on"
            @click="refresh"
          >
            <v-icon>refresh</v-icon>
          </v-btn>
        </template>
        <span>{{ $t('modals.service.refreshEntities') }}</span>
      </v-tooltip>
      <v-btn
        depressed
        text
        @click="$modals.hide"
      >
        {{ $t('common.close') }}
      </v-btn>
      <v-btn
        v-if="entitiesActionsInQueue"
        :loading="submitting"
        :disabled="submitting || !actionsRequests.length"
        class="primary"
        type="submit"
        @click="submit"
      >
        {{ $t('common.submit') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { PAGINATION_LIMIT } from '@/config';
import { MODALS, PBEHAVIOR_ORIGINS, USERS_PERMISSIONS } from '@/constants';

import Observer from '@/services/observer';

import { authMixin } from '@/mixins/auth';
import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesServiceEntityMixin } from '@/mixins/entities/service-entity';
import { entitiesAlarmTagMixin } from '@/mixins/entities/alarm-tag';
import { localQueryMixin } from '@/mixins/query/query';
import { permissionsWidgetsEventComment } from '@/mixins/permissions/widgets/entity-comment';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import ServiceTemplate from '@/components/other/service/partials/service-template.vue';
import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/pbehaviors-simple-list.vue';
import EntityComments from '@/components/other/entity/entity-comments.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.serviceEntities,
  provide() {
    return {
      $periodicRefresh: this.$periodicRefresh,
    };
  },
  inject: ['$system'],
  components: {
    PbehaviorsSimpleList,
    ServiceTemplate,
    EntityComments,
    ModalWrapper,
  },
  mixins: [
    authMixin,
    localQueryMixin,
    modalInnerMixin,
    entitiesAlarmTagMixin,
    entitiesServiceEntityMixin,
    permissionsWidgetsEventComment,
    submittableMixinCreator(),
    confirmableModalMixinCreator({ field: 'actionsRequests' }),
  ],
  data() {
    return {
      pending: true,
      unavailableEntitiesAction: {},
      actionsRequests: [],
      query: {
        itemsPerPage: this.modal.config.widgetParameters.modalItemsPerPage ?? PAGINATION_LIMIT,
        sortBy: ['state'],
        sortDesc: [true],
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

    entitiesActionsInQueue() {
      return this.widgetParameters?.entitiesActionsInQueue ?? false;
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
    refresh(immediate = false) {
      if (!immediate && this.entitiesActionsInQueue) {
        return Promise.resolve();
      }

      return this.$periodicRefresh.notify();
    },

    addAction(action) {
      this.actionsRequests.push(action);
    },

    clearActions(value) {
      if (!value) {
        this.actionsRequests = [];
      }
    },

    async fetchList() {
      this.pending = true;

      const params = this.getQuery();
      params.with_instructions = true;
      params.pbh_origin = PBEHAVIOR_ORIGINS.serviceWeather;

      await this.fetchServiceEntitiesList({ id: this.service._id, params });

      if (!this.alarmTagsPending) {
        this.fetchAlarmTagsList({ params: { paginate: false } });
      }

      this.pending = false;
    },

    async submit() {
      await Promise.all(this.actionsRequests.map(({ action }) => action()));

      this.$modals.hide();
    },
  },
};
</script>

<style lang="scss" scoped>
.actions-requests-alert {
  line-height: 15px;
}
</style>
