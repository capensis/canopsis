<template>
  <v-layout class="cluster-graph-editor">
    <v-flex class="cluster-graph-editor__sidebar">
      <cluster-graph-entities-type
        v-field="form.impact"
        class="mb-3"
        @change="clearPinnedEntities"
      />
      <cluster-graph-entities-list
        v-field="form.entities"
        :impact="form.impact"
        @remove="removeEntity"
        @update:entity="updateEntity"
        @update:pinned="updatePinnedEntities"
      />
    </v-flex>
    <v-flex class="cluster-graph-editor__content">
      <c-zoom-overlay>
        <network-graph
          ref="networkGraph"
          :options="options"
          :node-html-label-options="nodeHtmlLabelsOptions"
          class="fill-height"
          ctrl-wheel-zoom
        />
      </c-zoom-overlay>
    </v-flex>
  </v-layout>
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
      return this.form.entities.filter(({ entity }) => entity);
    },

    rootEntitiesById() {
      return keyBy(this.notEmptyEntities, 'entity._id');
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
      return this.notEmptyEntities.map(({ entity = {}, pinned = [] } = {}) => [
        entity._id,

        ...pinned.filter(({ _id: id }) => !this.rootEntitiesById[id]).map(({ _id: id }) => id),
      ]);
    },

    styleOption() {
      return [
        {
          selector: 'node',
          style: {
            width: TREE_OF_DEPENDENCIES_GRAPH_OPTIONS.nodeSize,
            height: TREE_OF_DEPENDENCIES_GRAPH_OPTIONS.nodeSize,
            'background-color': COLORS.secondary,
            'border-color': COLORS.secondary,
          },
        },
        {
          selector: 'node[showAll]',
          style: {
            'background-opacity': 0,
            'border-width': 0,
            width: 128,
            height: 34,
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
      return this.notEmptyEntities.reduce((acc, { entity, pinned = [] }) => {
        acc.push(
          {
            group: 'nodes',
            data: { id: entity._id, entity },
          },
        );

        pinned.forEach((pinnedEntity) => {
          acc.push(
            {
              group: 'nodes',
              data: { id: pinnedEntity._id, entity: pinnedEntity },
            },
            {
              group: 'edges',
              data: { source: entity._id, target: pinnedEntity._id },
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
     * @param {Object} entity
     */
    addRootElement(entity) {
      const elements = this.$refs.networkGraph.$cy.getElementById(entity._id);

      if (elements.length) {
        elements.data({ id: entity._id, entity });

        return;
      }

      this.$refs.networkGraph.$cy.add([
        {
          group: 'nodes',
          data: {
            id: entity._id,
            entity,
            root: true,
          },
        },
      ]);
    },

    /**
     * Update cluster root data in cytoscape
     *
     * @param {Object} entity
     */
    updateRootElement(entity) {
      this.$refs.networkGraph.$cy.getElementById(entity._id).data({
        id: entity._id,
        entity,
        root: true,
      });
    },

    /**
     * Remove cluster root from cytoscape
     *
     * @param {Object} entity
     * @param {Object[]} pinned
     */
    removeRootElement(entity, pinned) {
      this.removePinnedElements(pinned, entity._id);

      if (!this.allEntitiesById[entity._id]) {
        this.$refs.networkGraph.$cy.getElementById(entity._id).remove();
      }
    },

    /**
     * Remove special entity
     *
     * @param {Object} entity
     * @param {Object[]} pinned
     * @returns {Promise<void>}
     */
    async removeEntity({ entity, pinned }) {
      if (!entity) {
        return;
      }

      /**
       * @desc We've added nextTick to avoid problem with v-field usage
       */
      await this.$nextTick();

      this.removeRootElement(entity, pinned);
      this.runLayout();
    },

    /**
     * Update entity data for special entity
     *
     * @param {Object} [newEntityItem]
     * @param {Object} [oldEntityItem]
     * @returns {Promise<void>}
     */
    async updateEntity(newEntityItem, oldEntityItem) {
      /**
       * @desc We've added nextTick to avoid problem with v-field usage
       */
      await this.$nextTick();

      const { entity: oldEntity, pinned: oldPinned } = oldEntityItem;
      const { entity: newEntity } = newEntityItem;

      if (!newEntity) {
        if (oldEntity) {
          this.removeRootElement(oldEntity, oldPinned);
        }

        return;
      }

      if (oldEntity) {
        if (oldEntity._id !== newEntity._id) {
          this.removeRootElement(oldEntity, oldPinned);
        } else {
          this.updateRootElement(newEntity);
        }
      }

      this.addRootElement(newEntity);
      this.runLayout();
    },

    /**
     * Update pinned entities for special entity
     *
     * @param {Object} [newEntityItem]
     * @param {Object} [oldEntityItem]
     * @returns {Promise<void>}
     */
    async updatePinnedEntities(newEntityItem, oldEntityItem) {
      /**
       * @desc We've added nextTick to avoid problem with v-field usage
       */
      await this.$nextTick();

      const { entity, pinned: oldPinned } = oldEntityItem;
      const { pinned: newPinned } = newEntityItem;

      const newById = keyBy(newPinned, '_id');
      const oldById = keyBy(oldPinned, '_id');
      const added = newPinned.filter(({ _id: id }) => !oldById[id]);
      const removed = oldPinned.filter(({ _id: id }) => !newById[id]);

      this.addPinnedElements(added, entity._id);
      this.removePinnedElements(removed, entity._id);
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
