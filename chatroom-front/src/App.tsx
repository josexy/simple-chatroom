import { BrowserRouter, Route, Routes } from "react-router-dom";
import GContainer from "./components/GContainer";
import GIndex from "./components/GIndex";
import GHome from "./components/GHome";
import GRoomOutlet from "./components/GRoomOutlet";
import GChatRoom from "./components/GChatRoom";
import GNotFound from "./components/GNotFound";

import "./assets/css/style.css"

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<GContainer />}>
                    <Route index element={<GIndex />} />
                    <Route path="rooms" element={<GHome />} />
                    <Route path="room" element={<GRoomOutlet />}>
                        <Route path=":id" element={<GChatRoom />} />
                    </Route>
                    <Route path="*" element={<GNotFound />} />
                </Route>
            </Routes>
        </BrowserRouter>
    )
}

export default App;