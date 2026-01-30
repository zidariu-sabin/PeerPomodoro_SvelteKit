import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { type ViteDevServer, defineConfig } from 'vite';
import { playwright } from '@vitest/browser-playwright'

// import { Server } from 'socket.io'

// const webSocketServer = {
// 	name: 'webSocketServer',
// 	configureServer(server: ViteDevServer) {
// 		if (!server.httpServer) return

// 		const io = new Server(server.httpServer)

// 		io.on('connection', (socket) => {
// 			socket.emit('eventFromServer', 'Hello, World 👋')
// 		})
// 	}
// }

export default defineConfig({ plugins: [tailwindcss(), sveltekit(),]});
