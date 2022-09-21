<template lang="pug">
  div.tree-of-dependencies__preview
    c-zoom-overlay
      network-graph.fill-height(
        ref="networkGraph",
        :options="options",
        :node-html-label-options="nodeHtmlLabelsOptions",
        ctrl-wheel-zoom
      )
    c-help-icon.map-preview__help-icon(size="32", color="secondary", icon="help", top)
      div.pre-wrap(v-html="$t('flowchart.panzoom.helpText')")
</template>

<script>
import { omit } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { COLORS, PAGINATION_LIMIT } from '@/config';
import {
  MODALS,
  TREE_OF_DEPENDENCIES_GRAPH_OPTIONS,
  TREE_OF_DEPENDENCIES_GRAPH_LAYOUT_OPTIONS,
  TREE_OF_DEPENDENCIES_TYPES, ENTITY_TYPES,
} from '@/constants';

import { getEntityColor } from '@/helpers/color';
import { getTreeOfDependenciesEntityText, normalizeTreeOfDependenciesMapEntities } from '@/helpers/map';

// eslint-disable-next-line import/no-webpack-loader-syntax
import engineeringIcon from '!!svg-inline-loader?modules!@/assets/images/engineering.svg';

import NetworkGraph from '@/components/common/chart/network-graph.vue';

const { mapActions } = createNamespacedHelpers('service');

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
      entitiesById: normalizeTreeOfDependenciesMapEntities(this.map.parameters?.entities),
    };
  },
  computed: {
    impact() {
      return this.map.parameters?.type === TREE_OF_DEPENDENCIES_TYPES.impactChain;
    },

    rootEntities() {
      return Object.values(this.entitiesById)
        .filter(entity => entity.dependencies);
    },

    cytoscapeClusters() {
      return this.rootEntities.map(({ entity, dependencies }) => [
        entity._id,
        ...dependencies.filter(id => !this.entitiesById[id].dependencies?.length),
      ]);
    },

    entitiesElements() {
      return this.rootEntities.reduce((acc, { entity, dependencies = [] }) => {
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

        dependencies.forEach((childId) => {
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

        if (entity[this.countProperty] > dependencies.length) {
          acc.push(...this.getShowAllElements(entity));
        }

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
          'tree-of-dependencies__fetch-dependencies',
        );
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
        badgeEl.dataset.id = entity._id;

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

        if (this.hasDependencies(entity) && !root) {
          nodeEl.appendChild(getBadgeEl(entity, opened, pending));
        }

        return nodeEl.outerHTML;
      };

      const getShowAllContent = ({ entity }) => {
        const btnContentEl = document.createElement('div');
        btnContentEl.classList.add('v-btn__content');
        btnContentEl.textContent = `Show all (${entity[this.countProperty]})`; // TODO: Add translation

        const btnEl = document.createElement('button');
        btnEl.classList.add(
          'v-btn',
          'v-btn--round',
          'theme--light',
          'tree-of-dependencies__show-all-btn',
        );
        btnEl.appendChild(btnContentEl);

        return btnEl.outerHTML;
      };

      return [
        {
          query: 'node',
          valign: 'center',
          halign: 'center',
          tpl: getContent,
        },
        {
          query: 'node[showAll]',
          valign: 'center',
          halign: 'center',
          tpl: getShowAllContent,
        },
      ];
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
  },
  watch: {
    map(map) {
      this.entitiesById = normalizeTreeOfDependenciesMapEntities(map.parameters?.entities);

      setTimeout(() => { // TODO: change it
        this.resetLayout();
      }, 1000);
    },
  },
  mounted() {
    this.$refs.networkGraph.$cy.on('tap', this.tapHandler);
  },
  beforeDestroy() {
    this.$refs.networkGraph.$cy.off('tap', this.tapHandler);
  },
  methods: {
    ...mapActions({
      fetchServiceDependenciesWithoutStore: 'fetchDependenciesWithoutStore',
      fetchServiceImpactsWithoutStore: 'fetchImpactsWithoutStore',
    }),

    fetchDependenciesList(data) {
      return this.impact
        ? this.fetchServiceImpactsWithoutStore(data)
        : this.fetchServiceDependenciesWithoutStore(data);
    },

    hasDependencies(entity) {
      return !!entity[this.countProperty];
    },

    /**
     * Remove old elements and add new elements to network graph
     */
    resetLayout() {
      this.$refs.networkGraph.$cy.elements().remove();
      this.$refs.networkGraph.$cy.add(this.entitiesElements);
      this.runLayout();
    },

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

    /**
     * Get show all cytoscape elements for special entity
     *
     * @param {Entity} entity
     * @returns {[Object, Object]}
     */
    getShowAllElements(entity) {
      const showAllId = `show-all-${entity._id}`;

      return [
        {
          group: 'nodes',
          data: {
            id: showAllId,
            entity,
            showAll: true,
          },
        },
        {
          group: 'edges',
          data: {
            id: `show-all-edge-${entity._id}`,
            source: entity._id,
            target: showAllId,
          },
        },
      ];
    },

    /**
     * Add dependencies to cytoscape for special source
     *
     * @param {Entity[]} elements
     * @param {Entity} source
     */
    addDependenciesElements(elements, source) {
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
            source: source._id,
            target: element._id,
          },
        });

        return acc;
      }, []);

      if (source[this.countProperty] > elements.length) {
        addedElements.push(...this.getShowAllElements(source));
      }

      this.$refs.networkGraph.$cy.add(addedElements);
    },

    /**
     * Remove dependencies from cytoscape for special source
     *
     * @param {string[]} elementsIds
     * @param {string} sourceId
     */
    removeDependenciesElements(elementsIds, sourceId) {
      const nodesForRemoveSelectors = elementsIds.map(id => `node[id = "${id}"]`);
      nodesForRemoveSelectors.push(`node[id = "show-all-${sourceId}"]`);

      const nodesForRemove = this.$refs.networkGraph.$cy.elements(nodesForRemoveSelectors.join(','));
      const filteredNodesForRemove = nodesForRemove.filter(node => node.connectedEdges().size() === 1);

      const edgesForRemoveSelectors = elementsIds.map(id => `edge[source = "${sourceId}"][target = "${id}"]`);
      edgesForRemoveSelectors.push(`node[id = "show-all-edge-${sourceId}"]`);

      const edgesForRemove = this.$refs.networkGraph.$cy.elements(edgesForRemoveSelectors.join(','));

      filteredNodesForRemove.remove();
      edgesForRemove.remove();
    },

    /**
     * Method for dependencies fetching for special node
     *
     * @param {string} id
     */
    async fetchDependencies(id) {
      const target = this.$refs.networkGraph.$cy.getElementById(id);
      const { opened, entity, root, pending } = target.data();

      if (!this.hasDependencies(entity) || root || pending) {
        return;
      }

      if (opened) {
        const { dependencies } = this.entitiesById[id];

        this.removeDependenciesElements(dependencies, id);

        target.data({
          pending: false,
          opened: false,
        });

        this.$delete(this.entitiesById[id], 'dependencies');

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

      this.$set(this.entitiesById[id], 'dependencies', ids);

      this.addDependenciesElements(data, entity);
      this.runLayout();
    },

    /**
     * Show modal window with all entity dependencies
     *
     * @param {string} entityId
     */
    showAllDependenciesModal(entityId) {
      this.$modals.show({
        name: MODALS.entityDependenciesList,
        config: {
          entityId,
          impact: this.impact,
        },
      });
    },

    /**
     * Handler for tap event on whole cytoscape canvas
     *
     * @param {Object} target
     * @param {MouseEvent} originalEvent
     */
    tapHandler({ target, originalEvent }) {
      if (originalEvent.target.classList.contains('v-badge__badge')) {
        const { id } = originalEvent.target.dataset;

        if (id) {
          this.fetchDependencies(id);

          return;
        }
      }

      const { entity, showAll } = target.data();

      if (!showAll || !entity || !this.hasDependencies(entity)) {
        return;
      }

      this.showAllDependenciesModal(entity._id);
    },
  },
};
</script>
<style lang="scss" scoped>
.tree-of-dependencies__preview {
  position: relative;
  height: 800px;
  width: 100%;

  & /deep/ canvas[data-id='layer0-selectbox'] { // Hide selectbox layer from cytoscape
    display: none;
  }

  & /deep/ .v-badge__badge {
    top: -7px;
    right: -7px;

    * {
      pointer-events: none;
    }
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

  & /deep/ .tree-of-dependencies__fetch-dependencies {
    width: 100%;
    height: 100%;
    border-radius: 50%;
  }
}
</style>
