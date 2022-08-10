<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name")
    v-layout.mermaid-editor
      v-flex.mermaid-editor__sidebar
        div.mb-3
          h6.my-2.title.text-xs-center {{ $t('common.type') }}
          v-radio-group.justify-center(v-model="impact", row, hide-details)
            v-radio(:value="false", color="primary", label="Tree of dependencies")
            v-radio(:value="true", color="primary", label="Impact chain")
        div
          h6.my-2.title.text-xs-center {{ $tc('common.entity', 2) }}
          v-card.ma-2(v-for="(entity, index) in entities", :key="entity.key")
            v-card-text
              c-entity-field(
                :value="entity.data",
                :name="`entity-${entity.key}`",
                :entity-types="entityTypes",
                :item-disabled="isItemDisabled",
                item-text="name",
                return-object,
                clearable,
                @input="updateData($event, index)"
              )
              v-expand-transition
                v-combobox(
                  v-if="entity.data",
                  :value="entity.pinned",
                  :items="pinnedListByIds[entity.data._id]",
                  :loading="pendingByIds[entity.data._id]",
                  label="Pinned entities",
                  item-value="_id",
                  item-text="name",
                  deletable-chips,
                  chips,
                  multiple,
                  @change="updatePinnedEntity($event, index)"
                )
            v-card-actions
              v-layout(justify-end)
                c-action-btn(type="delete", @click="removeEntity(index)")
          v-btn(color="primary", @click="addEntity") Add entity
      v-flex.mermaid-editor__content
        v-layout.fill-height(column)
          network-graph.healthcheck-network-graph(
            ref="networkGraph",
            :options="options",
            :node-html-label-options="nodeHtmlLabelsOptions",
            @node:tap="nodeTapHandler",
            @node:mouseover="nodeMouseoverHandler",
            @node:mouseout="nodeMouseoutHandler"
          )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { ENTITY_TYPES, HEALTHCHECK_NETWORK_GRAPH_OPTIONS } from '@/constants';

import uid from '@/helpers/uid';

import { formMixin } from '@/mixins/form';

import MermaidEditor from '@/components/common/mermaid/mermaid-editor.vue';
import NetworkGraph from '@/components/common/chart/network-graph.vue';
import { PAGINATION_LIMIT } from '@/config';

const { mapActions } = createNamespacedHelpers('service');

// TODO: move to helpers
const toById = (entities = [], idKey = '_id') => entities.reduce((acc, entity) => {
  acc[entity[idKey]] = entity;

  return acc;
}, {});

export default {
  inject: ['$validator'],
  components: { MermaidEditor, NetworkGraph },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      impact: false,
      entities: [],
      pendingByIds: {},
      pinnedListByIds: {},
      pinnedListMetaByIds: {},
      query: {
        page: 1,
      },
    };
  },
  computed: {
    entityTypes() {
      return [ENTITY_TYPES.service];
    },

    styleOption() {
      return [
        {
          selector: 'node',
          style: {
            'font-size': '10px',
            'background-color': '#2B3E4F',
            'border-color': '#2B3E4F',
          },
        },
        {
          selector: 'edge',
          style: {
            width: 2,
            'curve-style': 'bezier',
            'line-color': 'black',
          },
        },
      ];
    },

    options() {
      return {
        style: this.styleOption,
        wheelSensitivity: HEALTHCHECK_NETWORK_GRAPH_OPTIONS.wheelSensitivity,
        minZoom: HEALTHCHECK_NETWORK_GRAPH_OPTIONS.minZoom,
        maxZoom: HEALTHCHECK_NETWORK_GRAPH_OPTIONS.maxZoom,
        layout: {
          // name: 'cise',
          // spacingFactor: HEALTHCHECK_NETWORK_GRAPH_OPTIONS.spacingFactor,
        },
      };
    },

    tooltipOptions() {
      return {};
    },

    nodeHtmlLabelsOptions() {
      return [
        {
          query: 'node',
          valign: 'top',
          halign: 'center',
          valignBox: 'top',
          halignBox: 'center',
          tpl: data => `<div>${data.entity.name}</div>`,
        },
      ];
    },

    allEntitiesById() {
      return this.entities.reduce((acc, { data = {}, pinned = [] } = {}) => {
        acc[data._id] = data;

        pinned.forEach((entity) => {
          acc[entity._id] = entity;
        });

        return acc;
      }, {});
    },
  },
  methods: {
    ...mapActions({
      fetchServiceDependenciesWithoutStore: 'fetchDependenciesWithoutStore',
      fetchServiceImpactsWithoutStore: 'fetchImpactsWithoutStore',
    }),

    isItemDisabled(item) {
      return !!this.entities.find(({ entity }) => entity._id === item._id);
    },

    addEntity() {
      this.entities.push({ key: uid(), entity: '', pinned: [] });
    },

    runLayout(fit = false) {
      if (this.entities.length === 0) {
        return;
      }

      this.$refs.networkGraph.$cy.layout({
        name: 'cise',
        animate: 'end',
        fit,
      }).run();
    },

    updateData(data, index) {
      const { data: oldData } = this.entities[index];

      this.$set(this.entities[index], 'data', data);
      this.$set(this.entities[index], 'pinned', []);

      if (!data) {
        if (oldData) {
          this.$refs.networkGraph.$cy.getElementById(oldData._id).remove();
        }

        return;
      }

      this.fetchPinnedEntitiesList(data._id);

      const elements = this.$refs.networkGraph.$cy.getElementById(data._id);

      if (elements.length) {
        elements.data({ id: data._id, entity: data });

        return;
      }

      this.$refs.networkGraph.$cy.add([
        {
          group: 'nodes',
          data: { id: data._id, entity: data },
        },
      ]);

      this.runLayout(true);
    },

    updatePinnedEntity(pinned, index) {
      const { data, pinned: oldPinned } = this.entities[index];

      this.$set(this.entities[index], 'pinned', pinned);

      const newById = toById(pinned);
      const oldById = toById(oldPinned);
      const added = pinned.filter(({ _id: id }) => !oldById[id]);
      const removed = oldPinned.filter(({ _id: id }) => !newById[id]);

      const addedItems = added.reduce((acc, item) => {
        const elements = this.$refs.networkGraph.$cy.getElementById(item._id);

        if (!elements.length) {
          acc.push({
            group: 'nodes',
            data: { id: item._id, entity: item },
          });
        }

        acc.push({
          group: 'edges',
          data: { source: data._id, target: item._id },
        });

        return acc;
      }, []);

      if (addedItems.length) {
        this.$refs.networkGraph.$cy.add(addedItems);
      }

      if (removed.length) {
        const removedNodes = removed.filter(({ _id: id }) => !this.allEntitiesById[id]);

        this.$refs.networkGraph.$cy.nodes().filter((el) => {
          const { id: elementId } = el.data();

          return removedNodes.find(({ _id: id }) => id === elementId);
        }).remove();

        this.$refs.networkGraph.$cy.edges().filter((el) => {
          const { source, target } = el.data();

          return removed.find(({ _id: id }) => id === target && data._id === source);
        }).remove();
      }

      this.runLayout();
    },

    removeEntity(index) {
      this.entities = this.entities.filter((item, ind) => index !== ind);
    },

    nodeTapHandler() {

    },

    nodeMouseoverHandler() {

    },

    nodeMouseoutHandler() {

    },

    async fetchPinnedEntitiesList(id, params = { limit: PAGINATION_LIMIT }) {
      this.$set(this.pendingByIds, id, true);

      const { data, meta } = await this.fetchDependenciesList({ id, params });

      this.$set(this.pinnedListByIds, id, data.map(({ entity }) => entity));
      this.$set(this.pinnedListMetaByIds, id, meta);
      this.$set(this.pendingByIds, id, false);
    },

    fetchDependenciesList(data) {
      return this.impact
        ? this.fetchServiceImpactsWithoutStore(data)
        : this.fetchServiceDependenciesWithoutStore(data);
    },
  },
};
</script>
