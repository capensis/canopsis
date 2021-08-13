import { API_SOCKET_HOST } from '@/config';

import { Socket } from '@/helpers/socket';

export default Socket.create({ baseUrl: API_SOCKET_HOST });
