import { useEffect, useRef, useState } from "react";
import * as PIXI from "pixi.js";

type AgentUpdate = {
  id: string;
  x: number;
  z: number;
  rotation: number;
  type: string;
  path?: Array<{ x: number; z: number }>;
};

type ObstacleUpdate = {
  id: string;
  minX: number;
  minZ: number;
  maxX: number;
  maxZ: number;
  type: string;
};

function getRotationFromVelocity(vx: number, vz: number) {
  return Math.atan2(vz, vx);
}
function App() {
  const [tick, setTick] = useState<number | null>(null);
  const stageRef = useRef<HTMLDivElement | null>(null);
  const appRef = useRef<PIXI.Application | null>(null);
  const container = useRef<PIXI.Container>(new PIXI.Container());
  const spritesRef = useRef<Map<string, PIXI.Graphics>>(new Map());
  const obstaclesRef = useRef<Map<string, PIXI.Graphics>>(new Map());
  const linesRef = useRef<Map<string, PIXI.Graphics>>(new Map());
  const targetsRef = useRef<Map<string, PIXI.Graphics>>(new Map());
  useEffect(() => {
    if (!stageRef.current) return;
    const app = new PIXI.Application();

    app
      .init({
        backgroundColor: 0x1e1e1e,
        resizeTo: stageRef.current,
        antialias: true,
        resolution: window.devicePixelRatio || 1,
      })
      .then(() => {
        if (!stageRef.current) return;
        stageRef.current.appendChild(app.canvas as HTMLCanvasElement);
        appRef.current = app;
        container.current.sortableChildren = true;
        app.stage.sortableChildren = true;
        app.stage.addChild(container.current);
      });

    return () => {
      app.destroy(true, { children: true });
      appRef.current = null;
    };
  }, []);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");
    ws.addEventListener("message", (ev) => {
      try {
        const data = JSON.parse(ev.data) as {
          tick: number;
          updated: AgentUpdate[];
          obstacles: ObstacleUpdate[];
        };
        setTick(data.tick);
        const app = appRef.current;
        if (!app) return;
        const W = app.renderer.width;
        const H = app.renderer.height;

        data.updated.forEach((u) => {
          if (u.type !== "agent") return;
          let g = spritesRef.current.get(u.id);
          const screenX = (u.x / 50) * W;
          const screenY = (u.z / 50) * H;

          if (!g) {
            g = new PIXI.Graphics();
            g.fill(0xffcc00);
            // g.drawRect(0, 0, 10, 10);
            const x1 = 2.5 * 2;
            const y1 = 2.5 * 2;
            const x2 = 1.25 * 2;
            const y2 = 5 * 2;
            const x3 = 3.75 * 2;
            const y3 = 5 * 2;

            g.moveTo(x1, y1);
            g.lineTo(x2, y2);
            g.lineTo(x3, y3);
            g.lineTo(x1, y1);
            g.fill();
            g.zIndex = 10;

            g.pivot.set(2.5 * 2, 3.75 * 2);

            g.x = screenX;
            g.y = screenY;

            g.rotation = u.rotation;

            container.current.addChild(g);
            spritesRef.current.set(u.id, g);
          } else {
            g.rotation = u.rotation;
            g.x += (screenX - g.x) * 0.6;
            g.y += (screenY - g.y) * 0.6;
          }

          // draw debug path waypoints
          let line = linesRef.current.get(u.id);
          if (!line) {
            line = new PIXI.Graphics();
            line.zIndex = 5;
            container.current.addChild(line);
            linesRef.current.set(u.id, line);
          }
          line.clear();
          line.zIndex = 5;

          if (u.path && u.path.length > 0) {
            // draw lines between waypoints
            line.setStrokeStyle({
              width: 2,
              color: 0x00ff00,
              alpha: 0.8,
            });
            let lastX = g.x;
            let lastY = g.y;
            for (let i = 0; i < u.path.length; i++) {
              const wp = u.path[i];
              const wpX = (wp.x / 50) * W;
              const wpY = (wp.z / 50) * H;
              line.moveTo(lastX, lastY).lineTo(wpX, wpY);
              lastX = wpX;
              lastY = wpY;
            }

            // draw waypoint markers
            let marker = targetsRef.current.get(u.id);
            if (!marker) {
              marker = new PIXI.Graphics();
              marker.zIndex = 4;
              container.current.addChild(marker);
              targetsRef.current.set(u.id, marker);
            }
            marker.clear();
            // draw small circles at each waypoint
            for (let i = 0; i < u.path.length; i++) {
              const wp = u.path[i];
              const wpX = (wp.x / 50) * W;
              const wpY = (wp.z / 50) * H;
              marker.fill(i === 0 ? 0xffff00 : 0x00ff00).circle(wpX, wpY, 3);
            }
            // final target in red
            if (u.path.length > 0) {
              const last = u.path[u.path.length - 1];
              const lastX = (last.x / 50) * W;
              const lastY = (last.z / 50) * H;
              marker.fill(0xff0000).circle(lastX, lastY, 4);
            }
          } else {
            // no path, clear markers
            const marker = targetsRef.current.get(u.id);
            if (marker) {
              marker.clear();
            }
          }
        });

        data.obstacles.forEach((o) => {
          if (o.type !== "obstacle") return;
          let g = obstaclesRef.current.get(o.id);
          if (!g) {
            g = new PIXI.Graphics();
            g.fill(0x0000ff); // blue color
            const width = Math.abs((o.maxX - o.minX) / 50) * W;
            const height = Math.abs((o.maxZ - o.minZ) / 50) * H;
            const screenX = (Math.min(o.minX, o.maxX) / 50) * W;
            const screenY = (Math.min(o.minZ, o.maxZ) / 50) * H;
            g.rect(0, 0, width, height);
            g.fill();
            g.x = screenX;
            g.y = screenY;
            g.tint = 0x0000ff;
            g.zIndex = 1; // behind agents
            container.current.addChild(g);
            obstaclesRef.current.set(o.id, g);
          }
          // obstacles don't move, so no update needed
        });
      } catch (e) {
        // ignore
      }
    });
    return () => ws.close();
  }, []);

  return (
    <div style={{ display: "flex", gap: 12 }}>
      <div style={{ width: 1000, height: 1000 }} ref={stageRef} />
      <div style={{ width: 220 }}>
        <div style={{ marginBottom: 8 }}>
          <strong>Tick:</strong> {tick ?? "-"}
        </div>
        <div>
          <em>Agents rendered with Pixi.js</em>
        </div>
      </div>
    </div>
  );
}

export default App;
