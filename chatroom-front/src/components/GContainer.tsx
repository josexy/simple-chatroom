import { Col, Container, Row } from "react-bootstrap";
import { Outlet } from "react-router-dom";

export default function GContainer() {
    return (
        <Container fluid className="pt-5 pb-5">
            <Row>
                <Col>
                    <Outlet />
                </Col>
            </Row>
        </Container>
    );
}
