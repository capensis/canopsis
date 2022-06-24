<template lang="pug">
  v-app#app
    v-layout(v-if="!pending")
      the-navigation#main-navigation(v-if="shownNavigation")
      v-content#main-content
        v-layout(row, wrap)
          v-btn.primary(@click="addRectangle") Add rect
          v-btn.primary(@click="addRoundedRectangle") Add rounded rect
          v-btn.primary(@click="addSquare") Add square
          v-divider(vertical)
          v-btn.primary(@click="addRhombus") Add rhombus
          v-btn.primary(@click="addParallelogram") Add parallelogram
          v-divider(vertical)
          v-btn.primary(@click="addCircle") Add circle
          v-btn.primary(@click="addEllipse") Add oval
          v-divider(vertical)
          v-btn.primary(@click="addLine") Add line
          v-btn.primary(@click="addArrowLine") Add arrow line
          v-btn.primary(@click="addBidirectionalArrowLine") Add bidirectional arrow line
          v-divider(vertical)
          v-btn.primary(@click="addText") Add text
          v-btn.primary(@click="addTextbox") Add textbox
          v-divider(vertical)
          v-btn.primary(@click="addStorage") Add storage
          v-divider(vertical)
          file-selector(
            ref="fileSelector",
            hide-details,
            @change="addImage"
          )
            template(#activator="{ on }")
              v-btn.primary(v-on="on") Add Image
        flowchart-editor(v-model="shapes")
    the-sidebar
    the-modals
    the-popups
</template>

<script>
import { isEmpty } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { SOCKET_URL, LOCAL_STORAGE_ACCESS_TOKEN_KEY } from '@/config';
import { EXCLUDED_SERVER_ERROR_STATUSES, MAX_LIMIT, ROUTES_NAMES, SHAPES } from '@/constants';

import { reloadPageWithTrailingSlashes } from '@/helpers/url';
import { convertDateToString } from '@/helpers/date/date';

import localStorageService from '@/services/local-storage';

import Socket from '@/plugins/socket/services/socket';

import { authMixin } from '@/mixins/auth';
import { systemMixin } from '@/mixins/system';
import { entitiesInfoMixin } from '@/mixins/entities/info';
import { entitiesUserMixin } from '@/mixins/entities/user';

import TheNavigation from '@/components/layout/navigation/the-navigation.vue';
import ActiveBroadcastMessage from '@/components/layout/broadcast-message/active-broadcast-message.vue';

import '@/assets/styles/main.scss';
import FlowchartEditor from '@/components/common/flowchart/c-flowchart-editor.vue';
import FileSelector from '@/components/forms/fields/file-selector.vue';
import { getFileDataUrlContent } from '@/helpers/file/file-select';
import { generatePoint } from '@/helpers/flowchart/points';

const { mapActions } = createNamespacedHelpers('remediationInstructionExecution');

export default {
  components: {
    FileSelector,
    FlowchartEditor,
    TheNavigation,
    ActiveBroadcastMessage,
  },
  mixins: [
    authMixin,
    systemMixin,
    entitiesInfoMixin,
    entitiesUserMixin,
  ],
  data() {
    return {
      shapes: {},
      pending: true,
    };
  },
  computed: {
    routeViewKey() {
      if (this.$route.name === ROUTES_NAMES.view) {
        return this.$route.path;
      }

      return this.$route.fullPath;
    },

    shownNavigation() {
      return ![ROUTES_NAMES.login, ROUTES_NAMES.error].includes(this.$route.name);
    },
  },
  beforeCreate() {
    reloadPageWithTrailingSlashes();
  },
  created() {
    this.registerCurrentUserOnceWatcher();
  },
  mounted() {
    this.socketConnectWithErrorHandling();
    this.fetchCurrentUserWithErrorHandling();
    this.showLocalStorageWarningPopupMessage();
  },
  methods: {
    ...mapActions({
      fetchPausedExecutionsWithoutStore: 'fetchPausedListWithoutStore',
    }),

    addRectangle() {
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.rect,
        width: 100,
        height: 100,
        x: 0,
        y: 0,
        text: '',
        alignCenter: true,
        justifyCenter: true,
        connections: [],
        style: {
          stroke: 'black',
          'stroke-width': 1,
          fill: 'white',
        },
      });
    },

    addRoundedRectangle() {
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.rect,
        width: 100,
        height: 100,
        x: 0,
        y: 0,
        rx: 20,
        ry: 20,
        text: '',
        alignCenter: true,
        justifyCenter: true,
        connections: [],
        style: {
          stroke: 'black',
          'stroke-width': 1,
          fill: 'white',
        },
      });
    },

    addLine() {
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.line,
        points: [
          generatePoint({
            x: 50,
            y: 50,
          }),
          generatePoint({
            x: 50,
            y: 150,
          }),
        ],
        connectedTo: [],
        style: {
          stroke: 'black',
          'stroke-width': 1,
        },
      });
    },

    addArrowLine() {
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.arrowLine,
        points: [
          generatePoint({
            x: 50,
            y: 50,
          }),
          generatePoint({
            x: 50,
            y: 150,
          }),
        ],
        connectedTo: [],
        style: {
          stroke: 'black',
          'stroke-width': 1,
        },
      });
    },

    addBidirectionalArrowLine() {
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.bidirectionalArrowLine,
        points: [
          generatePoint({
            x: 50,
            y: 50,
          }),
          generatePoint({
            x: 50,
            y: 150,
          }),
        ],
        connectedTo: [],
        style: {
          stroke: 'black',
          'stroke-width': 1,
        },
      });
    },

    addCircle() {
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.circle,
        x: 50,
        y: 50,
        diameter: 100,
        text: '',
        alignCenter: true,
        justifyCenter: true,
        connections: [],
        style: {
          stroke: 'black',
          'stroke-width': 1,
          fill: 'white',
        },
      });
    },

    addEllipse() {
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.ellipse,
        x: 50,
        y: 50,
        width: 100,
        height: 100,
        text: '',
        alignCenter: true,
        justifyCenter: true,
        connections: [],
        style: {
          stroke: 'black',
          'stroke-width': 1,
          fill: 'white',
        },
      });
    },

    addRhombus() {
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.rhombus,
        x: 50,
        y: 50,
        width: 100,
        height: 100,
        text: '',
        alignCenter: true,
        justifyCenter: true,
        connections: [],
        style: {
          stroke: 'black',
          'stroke-width': 1,
          fill: 'white',
        },
      });
    },

    addParallelogram() {
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.parallelogram,
        x: 50,
        y: 50,
        width: 100,
        height: 100,
        offset: 50,
        text: '',
        alignCenter: true,
        justifyCenter: true,
        connections: [],
        style: {
          stroke: 'black',
          'stroke-width': 1,
          fill: 'white',
        },
      });
    },

    addSquare() {
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.square,
        x: 50,
        y: 50,
        size: 100,
        text: '',
        connections: [],
        style: {
          stroke: 'black',
          'stroke-width': 1,
          fill: 'white',
        },
      });
    },

    addText() {
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.rect,
        x: 50,
        y: 50,
        width: 100,
        height: 100,
        text: 'Text',
        alignCenter: true,
        justifyCenter: true,
        connections: [],
        style: {
          fill: 'transparent',
        },
      });
    },

    addTextbox() {
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.rect,
        x: 50,
        y: 50,
        width: 120,
        height: 100,
        text: '<h1>Heading</h1><p>Paragraph</p>',
        connections: [],
        style: {
          fill: 'transparent',
        },
      });
    },

    addStorage() {
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.storage,
        x: 50,
        y: 50,
        width: 120,
        height: 100,
        radius: 20,
        text: 'Storage',
        alignCenter: true,
        justifyCenter: true,
        connections: [],
        style: {
          stroke: 'black',
          'stroke-width': 1,
          fill: 'grey',
        },
      });
    },

    async addImage([file]) {
      const url = await getFileDataUrlContent(file);
      const id = Date.now();

      this.$set(this.shapes, id, {
        id,
        type: SHAPES.image,
        x: 50,
        y: 50,
        width: 120,
        height: 100,
        src: url,
        text: file.name,
        justifyCenter: true,
        connections: [],
        style: {},
      });
    },

    showLocalStorageWarningPopupMessage() {
      const text = localStorageService.pop('warningPopup');

      if (text) {
        this.$popups.warning({ text, autoClose: false });
      }
    },

    registerCurrentUserOnceWatcher() {
      const unwatch = this.$watch('currentUser', async (currentUser) => {
        if (!isEmpty(currentUser)) {
          this.$socket.authenticate(localStorageService.get(LOCAL_STORAGE_ACCESS_TOKEN_KEY));

          await Promise.all([
            this.fetchAppInfo(),
            this.filesAccess(),
          ]);

          this.setSystemData({
            timezone: this.timezone,
          });

          this.setTitle();
          this.showPausedExecutionsPopup();

          unwatch();
        }
      });
    },

    async showPausedExecutionsPopup() {
      const pausedExecutions = await this.fetchPausedExecutionsWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      if (!pausedExecutions || !pausedExecutions.length) {
        return;
      }

      pausedExecutions.forEach((execution = {}) => this.$popups.info({
        text: this.$t('remediationInstructionExecute.popups.wasPaused', {
          instructionName: execution.instruction_name,
          alarmName: execution.alarm_name,
          date: convertDateToString(execution.paused),
        }),
      }));
    },

    socketConnectWithErrorHandling() {
      try {
        this.$socket
          .connect(SOCKET_URL)
          .on('error', this.socketErrorHandler);
      } catch (err) {
        this.$popups.error({
          text: this.$t('errors.socketConnectionProblem'),
          autoClose: false,
        });

        console.error(err);
      }
    },

    socketErrorHandler({ message } = {}) {
      if (message) {
        this.$popups.error({ text: message });

        if (message === Socket.ERROR_MESSAGES.authenticationFailed) {
          localStorageService.set('warningPopup', this.$t('warnings.authTokenExpired'));
          this.logout();
        }
      }
    },

    async fetchCurrentUserWithErrorHandling() {
      try {
        this.pending = true;

        await this.fetchCurrentUser();
      } catch (err) {
        if (!EXCLUDED_SERVER_ERROR_STATUSES.includes(err.status)) {
          this.$router.push({ name: ROUTES_NAMES.error });
        }

        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>

<style lang="scss">
  #app {
    &.-fullscreen {
      width: 100%;

      #main-navigation {
        display: none;
      }

      #main-content {
        padding: 0 !important;
      }
    }
  }
</style>
