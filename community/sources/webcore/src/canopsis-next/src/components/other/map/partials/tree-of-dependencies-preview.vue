<template lang="pug">
  c-zoom-overlay
    network-graph.fill-height(
      ref="networkGraph",
      :options="options",
      :node-html-label-options="nodeHtmlLabelsOptions",
      ctrl-wheel-zoom,
      @node:tap="nodeTapHandler"
    )
</template>

<script>
import { omit, cloneDeep } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { COLORS, PAGINATION_LIMIT } from '@/config';
import {
  TREE_OF_DEPENDENCIES_GRAPH_OPTIONS,
  TREE_OF_DEPENDENCIES_GRAPH_LAYOUT_OPTIONS,
  TREE_OF_DEPENDENCIES_TYPES,
} from '@/constants';

import { getTreeOfDependenciesEntityText } from '@/helpers/map';

import NetworkGraph from '@/components/common/chart/network-graph.vue';

const { mapActions } = createNamespacedHelpers('service');

const prepareEntities = (entities = []) => entities.reduce((acc, { entity, pinned_entities: pinnedEntities }) => {
  const newEntity = {
    entity: cloneDeep(entity),
    children: [],
  };

  pinnedEntities.forEach((pinnedEntity) => {
    const { _id: id } = pinnedEntity;

    newEntity.children.push(id);

    if (!acc[id]) {
      acc[id] = {
        entity: cloneDeep(pinnedEntity),
      };
    }
  });

  acc[entity._id] = newEntity;

  return acc;
}, {});

export default {
  components: { NetworkGraph },
  props: {
    map: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      entitiesById: prepareEntities(this.map.parameters?.entities),
    };
  },
  computed: {
    impact() {
      return this.map?.parameters.type === TREE_OF_DEPENDENCIES_TYPES.impactChain;
    },

    rootEntities() {
      return Object.values(this.entitiesById)
        .filter(entity => entity.children?.length);
    },

    cytoscapeClusters() {
      return this.rootEntities.map(({ entity, children }) => [
        entity._id,
        ...children.filter(id => !this.entitiesById[id].children?.length),
      ]);
    },

    entitiesElements() {
      return this.rootEntities.reduce((acc, { entity, children = [] }) => {
        acc.push(
          {
            group: 'nodes',
            data: {
              id: entity._id,
              entity,
              root: true,
            },
          },
        );

        children.forEach((childId) => {
          const child = this.entitiesById[childId];

          acc.push(
            {
              group: 'nodes',
              data: {
                id: childId,
                entity: child.entity,
              },
            },
            {
              group: 'edges',
              data: {
                source: entity._id,
                target: childId,
              },
            },
          );
        });

        return acc;
      }, []);
    },

    styleOption() {
      return [
        {
          selector: 'node',
          style: {
            width: TREE_OF_DEPENDENCIES_GRAPH_OPTIONS.nodeSize,
            height: TREE_OF_DEPENDENCIES_GRAPH_OPTIONS.nodeSize,
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

    countProperty() {
      return this.impact ? 'impacts_count' : 'depends_count';
    },

    nodeHtmlLabelsOptions() {
      const { nodeSize } = TREE_OF_DEPENDENCIES_GRAPH_OPTIONS;
      const getContent = ({ entity, opened, pending, root }) => {
        if (pending) {
          return '<div class="v-progress-circular v-progress-circular--indeterminate white--text position-relative" style="height: 100%; width: 100%;"><svg xmlns="http://www.w3.org/2000/svg" viewBox="22.857142857142858 22.857142857142858 45.714285714285715 45.714285714285715" style="transform: rotate(0deg);"><circle fill="transparent" cx="45.714285714285715" cy="45.714285714285715" r="15" stroke-width="3" stroke-dasharray="125.664" stroke-dashoffset="125.66370614359172px" class="v-progress-circular__overlay"></circle></svg><div class="v-progress-circular__info"></div></div>';
        }

        const count = entity[this.countProperty];

        return (count && !root)
          ? `<i class="v-icon material-icons theme--light white--text">${opened ? 'remove' : 'add'}</i>`
          : '';
      };

      return [
        {
          query: 'node',
          valign: 'center',
          halign: 'center',
          tpl: data => `<div class="secondary v-btn__content" style="width: ${nodeSize}px; height: ${nodeSize}px; border-radius: 50%;">${getContent(data)}<div class="position-absolute" style="top: ${nodeSize}px">${getTreeOfDependenciesEntityText(data.entity)}</div></div>`,
        },
      ];
    },

    options() {
      return {
        ...omit(TREE_OF_DEPENDENCIES_GRAPH_OPTIONS, ['nodeSize']),

        layout: {
          ...TREE_OF_DEPENDENCIES_GRAPH_LAYOUT_OPTIONS,

          clusters: this.cytoscapeClusters,
        },

        style: this.styleOption,
        elements: this.entitiesElements,
      };
    },
  },
  methods: {
    ...mapActions({
      fetchServiceDependenciesWithoutStore: 'fetchDependenciesWithoutStore',
      fetchServiceImpactsWithoutStore: 'fetchImpactsWithoutStore',
    }),

    /**
     * Run 'cise' layout for rerender clusters
     */
    runLayout() {
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

    addChildrenElements(elements, sourceId) {
      if (!elements.length) {
        return;
      }

      const addedElements = elements.reduce((acc, element) => {
        const items = this.$refs.networkGraph.$cy.getElementById(element._id);

        if (!items.length) {
          acc.push({
            group: 'nodes',
            data: {
              id: element._id,
              entity: element,
            },
          });
        }

        acc.push({
          group: 'edges',
          data: {
            source: sourceId,
            target: element._id,
          },
        });

        return acc;
      }, []);

      this.$refs.networkGraph.$cy.add(addedElements);
    },

    removeChildrenElements(elementsIds, sourceId) {
      const nodesForRemoveSelectors = elementsIds.map(id => `node[id = "${id}"]`);
      const nodesForRemove = this.$refs.networkGraph.$cy.elements(nodesForRemoveSelectors.join(','));
      const filteredNodesForRemove = nodesForRemove.filter(node => node.connectedEdges().size() === 1);

      const edgesForRemoveSelectors = elementsIds.map(id => `edge[source = "${sourceId}"][target = "${id}"]`);
      const edgesForRemove = this.$refs.networkGraph.$cy.elements(edgesForRemoveSelectors.join(','));

      filteredNodesForRemove.remove();
      edgesForRemove.remove();
    },

    /**
     * Handler for the 'tap' event on node for cytoscape
     *
     * @param {Object} target
     */
    async nodeTapHandler({ target }) {
      const { id, opened, entity, root } = target.data();

      if (!entity[this.countProperty] || root) {
        return;
      }

      if (opened) {
        const { children } = this.entitiesById[id];

        this.removeChildrenElements(children, id);

        target.data({
          pending: false,
          opened: false,
        });

        this.$delete(this.entitiesById[id], 'children');

        this.runLayout();

        return;
      }

      target.data({
        pending: true,
      });

      const { data } = await this.fetchDependenciesList({ id, params: { limit: PAGINATION_LIMIT } });

      target.data({
        opened: true,
        pending: false,
      });

      const ids = data.map((item) => {
        let newEntity = item;

        if (this.entitiesById[item._id]) {
          newEntity = { ...this.entitiesById[item._id], ...newEntity };
        }

        this.$set(this.entitiesById, item._id, newEntity);

        return item._id;
      });

      this.$set(this.entitiesById[id], 'children', ids);

      this.addChildrenElements(data, id);
      this.runLayout();
    },

    fetchDependenciesList(data) {
      return this.impact
        ? this.fetchServiceImpactsWithoutStore(data)
        : this.fetchServiceDependenciesWithoutStore(data);
    },
  },
};
</script>
