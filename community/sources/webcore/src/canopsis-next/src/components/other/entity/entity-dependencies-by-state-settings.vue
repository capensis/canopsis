<template>
  <div class="entity-dependencies-by-state-settings">
    <network-graph
      ref="networkGraph"
      :options="options"
      :node-html-label-options="nodeHtmlLabelsOptions"
      class="fill-height black--text"
      ctrl-wheel-zoom
    />
  </div>
</template>

<script>
import { omit } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import {
  ROOT_CAUSE_DIAGRAM_OPTIONS,
  ROOT_CAUSE_DIAGRAM_LAYOUT_OPTIONS,
  ENTITY_TYPES,
  ROOT_CAUSE_DIAGRAM_NODE_SIZE,
  ROOT_CAUSE_DIAGRAM_EVENTS_NODE_SIZE,
} from '@/constants';

import { getEntityColor } from '@/helpers/entities/entity/color';
import { getMapEntityText, normalizeTreeOfDependenciesMapEntities } from '@/helpers/entities/map/list';
import { isEntityEventsStateSettings } from '@/helpers/entities/entity/entity';
import { convertSortToRequest } from '@/helpers/entities/shared/query';

import { entitiesEntityDependenciesMixin } from '@/mixins/entities/entity-dependencies';

import NetworkGraph from '@/components/common/chart/network-graph.vue';

// eslint-disable-next-line import/no-webpack-loader-syntax
import engineeringIcon from '!!svg-inline-loader?modules!@/assets/images/engineering.svg';

const getIconElement = (node) => {
  const { isEvents, entity } = node;

  const el = document.createElement('i');
  el.classList.add(
    'v-icon',
    'material-icons',
    'theme--light',
    'white--text',
    'entity-dependencies-by-state-settings__node-icon',
  );

  if (isEvents) {
    el.classList.add('entity-dependencies-by-state-settings__node-icon--events');
  }

  el.innerHTML = isEvents
    ? 'textsms'
    : {
      [ENTITY_TYPES.service]: engineeringIcon,
      [ENTITY_TYPES.resource]: 'perm_identity',
      [ENTITY_TYPES.component]: 'developer_board',
    }[entity.type];

  return el;
};

const getProgressElement = () => {
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
    'v-progress-circular--visible',
    'white--text',
    'entity-dependencies-by-state-settings__node-progress',
  );

  return progressEl;
};

const getContentElement = (node) => {
  const { entity, pending, isEvents } = node;
  const nodeSize = isEvents ? ROOT_CAUSE_DIAGRAM_EVENTS_NODE_SIZE : ROOT_CAUSE_DIAGRAM_NODE_SIZE;

  const nodeLabelEl = document.createElement('div');
  nodeLabelEl.classList.add('position-absolute');
  nodeLabelEl.style.top = `${nodeSize}px`;

  if (!isEvents) {
    nodeLabelEl.textContent = getMapEntityText(entity);
  }

  const nodeEl = document.createElement('div');
  nodeEl.appendChild(getIconElement(node));
  nodeEl.appendChild(nodeLabelEl);
  nodeEl.classList.add('v-btn__content', 'position-relative', 'border-radius-rounded');
  nodeEl.style.width = `${nodeSize}px`;
  nodeEl.style.height = `${nodeSize}px`;
  nodeEl.style.justifyContent = 'center';
  nodeEl.style.background = getEntityColor(entity);

  if (pending) {
    nodeEl.appendChild(getProgressElement());
  }

  return nodeEl.outerHTML;
};

export default {
  components: { NetworkGraph },
  mixins: [entitiesEntityDependenciesMixin],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    colorIndicator: {
      type: String,
      required: false,
    },
    impact: {
      type: Boolean,
      required: false,
    },
    columns: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      metaByEntityId: {},
      entitiesById: normalizeTreeOfDependenciesMapEntities([{ entity: this.entity, pinned_entities: [] }]),
    };
  },
  computed: {
    isEventsStateSettings() {
      return isEntityEventsStateSettings(this.entity);
    },

    entitiesWithDependencies() {
      return Object.values(this.entitiesById).filter(entity => entity.dependencies);
    },

    entitiesElements() {
      return this.entitiesWithDependencies.reduce((acc, { entity, dependencies = [] }) => {
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

        if (isEntityEventsStateSettings(entity)) {
          acc.push(...this.getEventsNodeElementByEntity(entity));

          return acc;
        }

        if (dependencies.length) {
          acc.push(...this.getEntityDependenciesElement(entity, dependencies));
        }

        const meta = this.metaByEntityId[entity._id] ?? {};

        if (meta.page < meta.page_count) {
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
            width: ROOT_CAUSE_DIAGRAM_OPTIONS.nodeSize,
            height: ROOT_CAUSE_DIAGRAM_OPTIONS.nodeSize,
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
          selector: 'node[isEvents]',
          style: {
            width: 30,
            height: 30,
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

    nodeHtmlLabelsOptions() {
      const getShowAllContent = ({ entity }) => {
        const meta = this.metaByEntityId[entity._id] ?? {};

        const fetchedEntities = meta.page * meta.per_page;

        const btnContentEl = document.createElement('div');
        btnContentEl.classList.add('v-btn__content');
        btnContentEl.textContent = `Show more (${fetchedEntities} of ${meta.total_count})`;

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
          tpl: getContentElement,
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
        ...omit(ROOT_CAUSE_DIAGRAM_OPTIONS, ['nodeSize']),

        style: this.styleOption,
        elements: this.entitiesElements,
      };

      if (this.entitiesElements.length) {
        options.layout = {
          ...ROOT_CAUSE_DIAGRAM_LAYOUT_OPTIONS,
        };
      }

      return options;
    },
  },
  watch: {
    entity() {
      this.entitiesById = normalizeTreeOfDependenciesMapEntities([{ entity: this.entity, pinned_entities: [] }]);

      /**
       * TODO: investigate this behavior in the future
       */
      setTimeout(() => this.resetLayout(), 1000);
    },
  },
  async mounted() {
    this.$refs.networkGraph.$cy.on('tap', this.tapHandler);

    if (!this.isEventsStateSettings) {
      await this.fetchDependencies(this.entity._id);
    }

    this.runLayout();
  },
  beforeDestroy() {
    this.$refs.networkGraph.$cy.off('tap', this.tapHandler);
  },
  methods: {
    getEventsNodeElementByEntity(entity) {
      const eventsNodeId = `${entity._id}_events-node`;

      return [
        {
          group: 'nodes',
          data: {
            entity,
            id: eventsNodeId,
            isEvents: true,
          },
        },
        {
          group: 'edges',
          data: {
            source: entity._id,
            target: eventsNodeId,
          },
        },
      ];
    },

    getEntityDependenciesElement(entity, dependenciesIds = []) {
      const dependencies = dependenciesIds.map(id => this.entitiesById[id].entity);

      return dependencies.reduce((acc, child) => {
        acc.push(
          {
            group: 'nodes',
            data: {
              id: child._id,
              entity: child,
            },
          },
          {
            group: 'edges',
            data: {
              source: entity._id,
              target: child._id,
            },
          },
        );

        return acc;
      }, []);
    },

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
    async runLayout() {
      if (this.$refs.networkGraph.$cy.nodes().empty()) {
        return;
      }

      try {
        this.$refs.networkGraph.$cy.layout({
          ...ROOT_CAUSE_DIAGRAM_LAYOUT_OPTIONS,
        }).run();
      } catch (err) {
        console.warn(err);
      }
    },

    /**
     * Show dependencies for node
     *
     * @param {Object} target
     * @returns {Promise<void>}
     */
    async showDependencies(target) {
      const { id } = target.data();
      const { page } = this.metaByEntityId[id] ?? {};
      const newPage = page ? page + 1 : 1;

      target.data({
        pending: true,
      });

      const { data, meta } = await this.fetchDependenciesList({
        id,
        params: {
          page: newPage,
          limit: PAGINATION_LIMIT,
          with_flags: true,
          /**
           * TODO: Api doesn't support multi sort
           */
          ...convertSortToRequest(['last_update_date', 'state']),
        },
      });

      target.data({
        pending: false,
      });

      this.$set(this.metaByEntityId, id, meta);

      const ids = data.map((item) => {
        let newEntityItem = { entity: item };

        if (this.entitiesById[item._id]) {
          newEntityItem = {
            ...this.entitiesById[item._id],

            entity: {
              ...newEntityItem,
              ...this.entitiesById[item._id].entity,
            },
          };
        }

        this.$set(this.entitiesById, item._id, newEntityItem);

        return item._id;
      });

      this.$set(this.entitiesById[id], 'dependencies', [
        ...this.entitiesById[id].dependencies,
        ...ids,
      ]);

      this.$refs.networkGraph.$cy.elements('*').remove();
      this.$refs.networkGraph.$cy.add(this.entitiesElements);
    },

    /**
     * Method for dependencies fetching for special node
     *
     * @param {string} id
     */
    async fetchDependencies(id) {
      const target = this.$refs.networkGraph.$cy.getElementById(id);
      const { pending } = target.data();

      if (pending) {
        return;
      }

      await this.showDependencies(target);

      this.runLayout();
    },

    /**
     * Handler for tap event on whole cytoscape canvas
     *
     * @param {Object} target
     * @param {MouseEvent} originalEvent
     */
    tapHandler({ target }) {
      const { entity, showAll } = target.data();

      if (!showAll || !entity) {
        return;
      }

      this.fetchDependencies(entity._id);
    },
  },
};
</script>

<style lang="scss">
.entity-dependencies-by-state-settings {
  position: relative;
  height: 800px;
  width: 100%;
  border-radius: 5px;
  background: white;

  &__node-progress {
    position: absolute;
    inset: 0;
  }

  &__node-icon {
    --node-size: 30px;

    font-size: var(--node-size) !important;

    svg {
      height: var(--node-size) !important;
    }

    &--events {
      --node-size: 20px;
    }
  }

  &__fetch-dependencies {
    width: 100%;
    height: 100%;
    border-radius: 50%;
  }

  canvas[data-id='layer0-selectbox'] { // Hide selectbox layer from cytoscape
    display: none;
  }

  .v-badge__badge {
    top: -7px;
    right: -7px;

    * {
      pointer-events: none;
    }
  }
}
</style>
