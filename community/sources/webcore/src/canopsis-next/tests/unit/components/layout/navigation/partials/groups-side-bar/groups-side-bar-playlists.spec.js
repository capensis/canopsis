import { generateRenderer } from '@unit/utils/vue';
import { createEntitiesModule, createMockedStoreModules, createPlaylistModule } from '@unit/utils/store';

import GroupsSideBarPlaylists from '@/components/layout/navigation/partials/groups-side-bar/groups-side-bar-playlists.vue';

const snapshotStubs = {
  'router-link': true,
};

describe('groups-side-bar-playlists', () => {
  const { playlistModule, playlists } = createPlaylistModule();
  const { entitiesModule } = createEntitiesModule();
  const store = createMockedStoreModules([entitiesModule, playlistModule]);

  const snapshotFactory = generateRenderer(GroupsSideBarPlaylists, { stubs: snapshotStubs });

  it('Renders `groups-side-bar-playlists` with required props', async () => {
    const wrapper = snapshotFactory({ store });

    await wrapper.openAllExpansionPanels();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `groups-side-bar-playlists` with custom props', async () => {
    playlists.mockReturnValue([
      {
        _id: 'playlist-1-id',
        name: 'playlist-1-name',
        enabled: true,
      },
      {
        _id: 'playlist-2-id',
        name: 'playlist-2-name',
        enabled: true,
      },
      {
        _id: 'playlist-3-id',
        name: 'playlist-3-name',
        enabled: false,
      },
    ]);
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([entitiesModule, playlistModule]),
    });

    await wrapper.openAllExpansionPanels();

    expect(wrapper).toMatchSnapshot();
  });
});
