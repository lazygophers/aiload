import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Home from './pages/Home';
import Docs from './pages/docs.tsx';
import Document from './pages/Document';
import './App.css';

function App() {
	return (
		<Router>
			<Layout>
				<Routes>
					<Route path="/" element={<Home />} />
					<Route path="/home" element={<Home />} />
					<Route path="/docs" element={<Docs />} />
					<Route path="/docs/:docId" element={<Document />} />
					{/* 可以在这里添加更多路由 */}
				</Routes>
			</Layout>
		</Router>
	);
}

export default App;
