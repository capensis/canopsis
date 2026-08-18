<template>
  <c-modal-wrapper
    :title-class="titleClass"
    close
  >
    <template #title>
      <span>{{ service.name }}</span>
    </template>
    <template #text>
      <v-tabs
        v-model="activeTab"
        align-tabs="center"
        color="primary"
        class="position-relative"
      >
        <v-tab :value="tabs.service">{{ $t('common.service') }}</v-tab>
        <v-tab :value="tabs.pbehavior" :disabled="!hasPbehaviorListAccess">
          {{ $t('common.pbehavior', 2) }}
        </v-tab>
        <v-tab v-if="hasAccessToComment" :value="tabs.comments">
          {{ $t('common.comment', 2) }}
        </v-tab>
      </v-tabs>

      <v-window v-model="activeTab" class="mt-2">
        <v-window-item :value="tabs.service">
          <c-progress-overlay :pending="pending" />
          <service-template
            v-model:options="options"
            :service="service"
            :service-entities="serviceEntitiesWithKey"
            :widget-parameters="widgetParameters"
            :total-items="serviceEntitiesMeta.total_count"
            :actions-requests="actionsRequests"
            @refresh="refresh"
            @add:action="addAction"
          />
        </v-window-item>

        <v-window-item :value="tabs.pbehavior">
          <pbehaviors-simple-list
            :entity="service"
            :default-inherited="widgetParameters.defaultInheritedPbehavior"
            with-active-status
            with-inherited
            addable
            updatable
            removable
            dense
          />
        </v-window-item>

        <v-window-item v-if="hasAccessToComment" :value="tabs.comments">
          <entity-comments
            :entity="service"
            addable
            editable
          />
        </v-window-item>
      </v-window>
    </template>
    <template #actions>
      <v-alert
        v-if="actionsRequests.length > 0"
        :model-value="actionsRequests.length > 0"
        class="actions-requests-alert my-0 mr-2 pa-1 pr-2"
        type="info"
        closable
        @update:model-value="clearActions"
      >
        {{ actionsRequests.length }} {{ $t('modals.service.actionInQueue', actionsRequests.length) }}
      </v-alert>
      <v-tooltip location="top">
        <template #activator="{ props: tooltipProps }">
          <v-btn
            class="mx-2"
            color="secondary"
            v-bind="tooltipProps"
            @click="refresh"
          >
            <v-icon icon="refresh" />
          </v-btn>
        </template>
        <span>{{ $t('modals.service.refreshEntities') }}</span>
      </v-tooltip>
      <v-btn
        variant="text"
        @click="close"
      >
        {{ $t('common.close') }}
      </v-btn>
      <v-btn
        v-if="entitiesActionsInQueue"
        :loading="submitting"
        :disabled="submitting || !actionsRequests.length"
        color="primary"
        type="submit"
        @click="submit"
      >
        {{ $t('common.submit') }}
      </v-btn>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { mapActions, mapState } from 'pinia';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS, PBEHAVIOR_ORIGINS, USER_PERMISSIONS } from '@/constants';

import Observer from '@/services/observer';

import { useServiceEntityStore } from '@/stores/modules/service-entity';

import { authMixin } from '@/mixins/auth';
import { permissionsWidgetsEventComment } from '@/mixins/permissions/widgets/entity-comment';
import { localQueryMixin } from '@/mixins/query/query';

import EntityComments from '@/components/other/entity/entity-comments.vue';
import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/pbehaviors-simple-list.vue';
import ServiceTemplate from '@/components/other/service/partials/service-template.vue';

const tabs = {
  service: 'service',
  pbehavior: 'pbehavior',
  comments: 'comments',
};

export default {
  name: MODALS.serviceEntities,

  inject: ['$system'],

  provide() {
    return {
      $periodicRefresh: this.$periodicRefresh,
    };
  },

  components: {
    PbehaviorsSimpleList,
    ServiceTemplate,
    EntityComments,
  },

  mixins: [
    authMixin,
    localQueryMixin,
    permissionsWidgetsEventComment,
  ],

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  data() {
    const { config = {} } = this.modal;

    return {
      pending: true,
      submitting: false,
      actionsRequests: [],
      activeTab: tabs.service,
      tabs,
      query: {
        page: 1,
        itemsPerPage: config.widgetParameters?.modalItemsPerPage ?? PAGINATION_LIMIT,
        sortBy: ['state'],
        sortDesc: [true],
      },
    };
  },

  computed: {
    ...mapState(useServiceEntityStore, {
      getServiceEntitiesListByServiceId: 'getListByServiceId',
      getServiceEntitiesMetaByServiceId: 'getMetaByServiceId',
    }),

    config() {
      return this.modal.config ?? {};
    },

    service() {
      return this.config.service;
    },

    titleClass() {
      return this.config.titleClass;
    },

    widgetParameters() {
      return this.config.widgetParameters ?? {};
    },

    entitiesActionsInQueue() {
      return this.widgetParameters?.entitiesActionsInQueue ?? false;
    },

    serviceEntities() {
      return this.getServiceEntitiesListByServiceId(this.service._id);
    },

    serviceEntitiesMeta() {
      return this.getServiceEntitiesMetaByServiceId(this.service._id);
    },

    serviceEntitiesWithKey() {
      return this.serviceEntities.map(entity => ({
        ...entity,
        key: `${entity._id}_${entity.alarm_id}`,
      }));
    },

    hasPbehaviorListAccess() {
      return this.checkAccess(USER_PERMISSIONS.business.serviceWeather.actions.pbehaviorList);
    },
  },

  beforeCreate() {
    this.$periodicRefresh = new Observer();
  },

  created() {
    this.$periodicRefresh.register(this.fetchList);
  },

  beforeUnmount() {
    this.$periodicRefresh.unregister(this.fetchList);
  },

  async mounted() {
    await this.fetchList();
  },

  methods: {
    ...mapActions(useServiceEntityStore, {
      fetchServiceEntitiesList: 'fetchList',
    }),

    close() {
      this.$modals.hide({ id: this.modal.id });
    },

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
      params.with_tag_colors = true;
      params.pbh_origin = PBEHAVIOR_ORIGINS.serviceWeather;

      await this.fetchServiceEntitiesList({ id: this.service._id, params });

      this.pending = false;
    },

    async submit() {
      this.submitting = true;

      try {
        await Promise.all(this.actionsRequests.map(({ action }) => action()));

        this.close();
      } finally {
        this.submitting = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.actions-requests-alert {
  line-height: 15px;
}
</style>
