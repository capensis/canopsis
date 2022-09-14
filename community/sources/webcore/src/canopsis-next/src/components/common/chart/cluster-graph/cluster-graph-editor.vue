<template lang="pug">
  v-layout.cluster-graph-editor
    v-flex.cluster-graph-editor__sidebar
      cluster-graph-entities-type.mb-3(
        v-field="form.impact",
        @change="clearPinnedEntities"
      )
      cluster-graph-entities-list(
        v-field="form.entities",
        :impact="form.impact",
        @remove="removeEntity",
        @update:data="updateEntityData",
        @update:pinned="updatePinnedEntities"
      )
    v-flex.cluster-graph-editor__content
      c-zoom-overlay
        network-graph.fill-height(
          ref="networkGraph",
          :options="options",
          :node-html-label-options="nodeHtmlLabelsOptions",
          ctrl-wheel-zoom
        )
</template>

<script>
import { keyBy, omit } from 'lodash';

import { COLORS } from '@/config';

import {
  TREE_OF_DEPENDENCIES_GRAPH_OPTIONS,
  TREE_OF_DEPENDENCIES_GRAPH_LAYOUT_OPTIONS,
  ENTITY_TYPES,
} from '@/constants';

import { formMixin } from '@/mixins/form';

import NetworkGraph from '@/components/common/chart/network-graph.vue';

import ClusterGraphEntitiesType from './partials/cluster-graph-entities-type.vue';
import ClusterGraphEntitiesList from './partials/cluster-graph-entities-list.vue';

export default {
  inject: ['$validator'],
  components: {
    NetworkGraph,
    ClusterGraphEntitiesType,
    ClusterGraphEntitiesList,
  },
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
  computed: {
    notEmptyEntities() {
      return this.form.entities.filter(({ data }) => data);
    },

    rootEntitiesById() {
      return keyBy(this.notEmptyEntities, 'data._id');
    },

    allPinnedEntitiesById() {
      return this.notEmptyEntities.reduce((acc, { pinned = [] } = {}) => {
        pinned.forEach((entity) => {
          acc[entity._id] = entity;
        });

        return acc;
      }, {});
    },

    allEntitiesById() {
      return { ...this.rootEntitiesById, ...this.allPinnedEntitiesById };
    },

    cytoscapeClusters() {
      return this.notEmptyEntities.map(({ data = {}, pinned = [] } = {}) => [
        data._id,

        ...pinned.filter(({ _id: id }) => !this.rootEntitiesById[id]).map(({ _id: id }) => id),
      ]);
    },

    styleOption() {
      return [
        {
          selector: 'node',
          style: {
            'font-size': '10px',
            'background-color': COLORS.secondary,
            'border-color': COLORS.secondary,
          },
        },
        {
          selector: 'edge',
          style: {
            width: 2,
            'curve-style': 'bezier',
            'line-color': 'silver',
          },
        },
      ];
    },

    entitiesElements() {
      return this.notEmptyEntities.reduce((acc, { data, pinned = [] }) => {
        acc.push(
          {
            group: 'nodes',
            data: { id: data._id, entity: data },
          },
        );

        pinned.forEach((entity) => {
          acc.push(
            {
              group: 'nodes',
              data: { id: entity._id, entity },
            },
            {
              group: 'edges',
              data: { source: data._id, target: entity._id },
            },
          );
        });

        return acc;
      }, []);
    },

    options() {
      const options = {
        ...omit(TREE_OF_DEPENDENCIES_GRAPH_OPTIONS, ['nodeSize']),

        style: this.styleOption,
        elements: this.entitiesElements,
      };

      if (this.entitiesElements.length) {
        options.layout = {
          ...TREE_OF_DEPENDENCIES_GRAPH_LAYOUT_OPTIONS,

          clusters: this.cytoscapeClusters,
        };
      }

      return options;
    },

    nodeHtmlLabelsOptions() {
      const tpl = ({ entity }) => `<div>${entity.type === ENTITY_TYPES.service ? entity.name : entity._id}</div>`;

      return [
        {
          tpl,
          query: 'node',
          valign: 'top',
          halign: 'center',
          valignBox: 'top',
          halignBox: 'center',
        },
      ];
    },
  },
  methods: {
    /**
     * Run 'cise' layout for rerender clusters
     */
    async runLayout() {
      await this.$nextTick();

      if (this.$refs.networkGraph.$cy.nodes().empty()) {
        return;
      }

      try {
        this.$refs.networkGraph.$cy.layout({
          ...TREE_OF_DEPENDENCIES_GRAPH_LAYOUT_OPTIONS,

          clusters: this.cytoscapeClusters,
        }).run();
      } catch (err) {
        console.warn(err);
      }
    },

    /**
     * Add pinned elements to cytoscape
     *
     * @param {Object[]} elements
     * @param {string} sourceId
     */
    addPinnedElements(elements = [], sourceId) {
      if (!elements.length) {
        return;
      }

      const addedElements = elements.reduce((acc, element) => {
        const items = this.$refs.networkGraph.$cy.getElementById(element._id);

        if (!items.length) {
          acc.push({
            group: 'nodes',
            data: { id: element._id, entity: element },
          });
        }

        acc.push({
          group: 'edges',
          data: { source: sourceId, target: element._id },
        });

        return acc;
      }, []);

      this.$refs.networkGraph.$cy.add(addedElements);
    },

    /**
     * Remove pinned elements from cytoscape
     *
     * @param {Object[]} elements
     * @param {string} sourceId
     */
    removePinnedElements(elements = [], sourceId) {
      if (!elements.length) {
        return;
      }

      const nodesForRemove = elements.filter(({ _id: id }) => !this.allEntitiesById[id])
        .map(({ _id: id }) => `node[id = "${id}"]`);

      const edgesForRemove = elements.map(({ _id: id }) => `edge[source = "${sourceId}"][target = "${id}"]`);
      const elementsForRemoveSelector = [...nodesForRemove, ...edgesForRemove].join(',');

      this.$refs.networkGraph.$cy.elements(elementsForRemoveSelector).remove();
    },

    /**
     * Add cluster root to cytoscape
     *
     * @param {Object} data
     */
    addRootElement(data) {
      const elements = this.$refs.networkGraph.$cy.getElementById(data._id);

      if (elements.length) {
        elements.data({ id: data._id, entity: data });

        return;
      }

      this.$refs.networkGraph.$cy.add([
        {
          group: 'nodes',
          data: {
            id: data._id,
            entity: data,
            root: true,
          },
        },
      ]);
    },

    /**
     * Update cluster root data in cytoscape
     *
     * @param {Object} data
     */
    updateRootElement(data) {
      this.$refs.networkGraph.$cy.getElementById(data._id).data({
        id: data._id,
        entity: data,
        root: true,
      });
    },

    /**
     * Remove cluster root from cytoscape
     *
     * @param {Object} data
     * @param {Object[]} pinned
     */
    removeRootElement(data, pinned) {
      this.removePinnedElements(pinned, data._id);

      if (!this.allEntitiesById[data._id]) {
        this.$refs.networkGraph.$cy.getElementById(data._id).remove();
      }
    },

    /**
     * Remove special entity
     *
     * @param {Object} data
     * @param {Object[]} pinned
     * @returns {Promise<void>}
     */
    async removeEntity({ data, pinned }) {
      if (!data) {
        return;
      }

      /**
       * @desc We've added nextTick to avoid problem with v-field usage
       */
      await this.$nextTick();

      this.removeRootElement(data, pinned);
      this.runLayout();
    },

    /**
     * Update entity data for special entity
     *
     * @param {Object} [newEntity]
     * @param {Object} [oldEntity]
     * @returns {Promise<void>}
     */
    async updateEntityData(newEntity, oldEntity) {
      /**
       * @desc We've added nextTick to avoid problem with v-field usage
       */
      await this.$nextTick();

      const { data: oldData, pinned: oldPinned } = oldEntity;
      const { data: newData } = newEntity;

      if (!newData) {
        if (oldData) {
          this.removeRootElement(oldData, oldPinned);
        }

        return;
      }

      if (oldData) {
        if (oldData._id !== newData._id) {
          this.removeRootElement(oldData, oldPinned);
        } else {
          this.updateRootElement(newData);
        }
      }

      this.addRootElement(newData);
      this.runLayout();
    },

    /**
     * Update pinned entities for special entity
     *
     * @param {Object} [newEntity]
     * @param {Object} [oldEntity]
     * @returns {Promise<void>}
     */
    async updatePinnedEntities(newEntity, oldEntity) {
      /**
       * @desc We've added nextTick to avoid problem with v-field usage
       */
      await this.$nextTick();

      const { data, pinned: oldPinned } = oldEntity;
      const { pinned: newPinned } = newEntity;

      const newById = keyBy(newPinned, '_id');
      const oldById = keyBy(oldPinned, '_id');
      const added = newPinned.filter(({ _id: id }) => !oldById[id]);
      const removed = oldPinned.filter(({ _id: id }) => !newById[id]);

      this.addPinnedElements(added, data._id);
      this.removePinnedElements(removed, data._id);
      this.runLayout();
    },

    /**
     * Clear pinned entities on change entities type
     *
     * @returns {Promise<void>}
     */
    async clearPinnedEntities() {
      /**
       * @desc We've added nextTick to avoid problem with v-field usage
       */
      await this.$nextTick();

      const newEntities = this.form.entities.map(entity => ({
        ...entity,

        pinned: [],
      }));

      this.updateField('entities', newEntities);

      this.$refs.networkGraph.$cy.elements('edge, node[!root]').remove();

      this.runLayout();
    },
  },
};
</script>

<style lang="scss" scoped>
$height: 500px;
$sideBarWidth: 500px;
$contentWidth: 800px;
$borderColor: #e5e5e5;

.cluster-graph-editor {
  &__sidebar, &__content {
    height: $height;
    border: 1px solid $borderColor;
  }

  &__sidebar {
    width: $sideBarWidth;
    overflow-y: auto;
    overflow-x: hidden;
    border-right: none;
  }

  &__content {
    width: $contentWidth;
  }
}
</style>
