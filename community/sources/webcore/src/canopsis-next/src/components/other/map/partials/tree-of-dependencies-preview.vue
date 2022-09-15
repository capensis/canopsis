<template lang="pug">
  c-zoom-overlay.tree-of-dependencies__preview
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
  TREE_OF_DEPENDENCIES_TYPES, ENTITY_TYPES,
} from '@/constants';

import { getTreeOfDependenciesEntityText } from '@/helpers/map';
import { getEntityColor } from '@/helpers/color';

// eslint-disable-next-line import/no-webpack-loader-syntax
import engineeringIcon from '!!svg-inline-loader?modules!@/assets/images/engineering.svg';

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
    colorIndicator: {
      type: String,
      required: false,
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

      const getIconEl = (entity) => {
        const el = document.createElement('i');
        el.classList.add(
          'v-icon',
          'material-icons',
          'theme--light',
          'white--text',
          'tree-of-dependencies__node-icon',
        );

        el.innerHTML = entity.type === ENTITY_TYPES.service
          ? engineeringIcon
          : 'person';

        return el;
      };

      const getProgressEl = () => {
        const progressContentCircleEl = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
        progressContentCircleEl.classList.add('v-progress-circular__overlay');
        progressContentCircleEl.setAttribute('fill', 'transparent');
        progressContentCircleEl.setAttribute('cx', '45.714285714285715');
        progressContentCircleEl.setAttribute('cy', '45.714285714285715');
        progressContentCircleEl.setAttribute('r', '15');
        progressContentCircleEl.setAttribute('stroke-width', '3');
        progressContentCircleEl.setAttribute('stroke-dasharray', '125.664');
        progressContentCircleEl.setAttribute('stroke-dashoffset', '125.66370614359172px');

        const progressContentEl = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
        progressContentEl.setAttribute('viewBox', '22.857142857142858 22.857142857142858 45.714285714285715 45.714285714285715');
        progressContentEl.appendChild(progressContentCircleEl);

        const progressEl = document.createElement('div');
        progressEl.appendChild(progressContentEl);
        progressEl.classList.add(
          'v-progress-circular',
          'v-progress-circular--indeterminate',
          'white--text',
          'position-relative',
        );

        return progressEl;
      };

      const getBadgeIconEl = (entity, opened) => {
        const badgeIconEl = document.createElement('i');
        badgeIconEl.classList.add(
          'v-icon',
          'material-icons',
          'theme--light',
          'white--text',
          'tree-of-dependencies__load-children',
        );
        badgeIconEl.dataset.id = entity._id;
        badgeIconEl.textContent = opened ? 'remove' : 'add';

        return badgeIconEl;
      };

      const getBadgeEl = (entity, opened, pending) => {
        const badgeEl = document.createElement('span');
        badgeEl.appendChild(pending ? getProgressEl() : getBadgeIconEl(entity, opened));
        badgeEl.classList.add(
          'v-badge__badge',
          'grey',
          'darken-1',
          'cursor-pointer',
        );

        return badgeEl;
      };

      const getContent = ({ entity, opened, pending, root }) => {
        const nodeLabelEl = document.createElement('div');
        nodeLabelEl.classList.add('position-absolute');
        nodeLabelEl.style.top = `${nodeSize}px`;
        nodeLabelEl.textContent = getTreeOfDependenciesEntityText(entity);

        const nodeEl = document.createElement('div');
        nodeEl.appendChild(getIconEl(entity));
        nodeEl.appendChild(nodeLabelEl);
        nodeEl.classList.add('v-btn__content', 'position-relative', 'border-radius-rounded');
        nodeEl.style.width = `${nodeSize}px`;
        nodeEl.style.height = `${nodeSize}px`;
        nodeEl.style.background = !this.colorIndicator
          ? COLORS.secondary
          : getEntityColor(entity, this.colorIndicator);

        if (entity[this.countProperty] && !root) {
          nodeEl.appendChild(getBadgeEl(entity, opened, pending));
        }

        return nodeEl.outerHTML;
      };

      return [
        {
          query: 'node',
          valign: 'center',
          halign: 'center',
          tpl: data => getContent(data),
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
  mounted() {
    this.$refs.networkGraph.$cy.on('mousedown', this.mousedownHandler);
  },
  beforeDestroy() {
    this.$refs.networkGraph.$cy.off('mousedown', this.mousedownHandler);
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
     * Method for children loading for special node
     *
     * @param {string} id
     */
    async loadChildrenById(id) {
      const target = this.$refs.networkGraph.$cy.elements(`node[id = "${id}"]`);

      const { opened, entity, root } = target.data();

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

    /**
     * Handler for cytoscape `mousedown` event
     *
     * @param {Object} event
     */
    mousedownHandler(event) {
      const { originalEvent } = event;
      if (!originalEvent.target.classList.contains('tree-of-dependencies__load-children')) {
        return;
      }

      originalEvent.preventDefault();

      this.loadChildrenById(originalEvent.target.dataset.id);
    },

    nodeTapHandler() {},

    fetchDependenciesList(data) {
      return this.impact
        ? this.fetchServiceImpactsWithoutStore(data)
        : this.fetchServiceDependenciesWithoutStore(data);
    },
  },
};
</script>
<style lang="scss" scoped>
.tree-of-dependencies__preview {
  & /deep/ canvas[data-id='layer0-selectbox'] { // Hide selectbox layer from cytoscape
    display: none;
  }

  & /deep/ .v-badge__badge {
    top: -7px;
    right: -7px;
  }

  & /deep/ .v-progress-circular {
    width: 20px;
    height: 20px;
  }

  & /deep/ .tree-of-dependencies__node-icon {
    font-size: 30px;

    svg {
      height: 30px;
    }
  }

  & /deep/ .tree-of-dependencies__load-children {
    width: 100%;
    height: 100%;
    border-radius: 50%;
  }
}
</style>
