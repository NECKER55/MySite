package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Configurazione della griglia Perlin
const (
	GRID_SLOTS       = 10 // Numero di slot della griglia
	PIXELS_PER_SLOT  = 10 // Pixel per ogni slot
	WORLD_HEIGHT     = 70 // Altezza del mondo
	WORLD_LENGTH     = 100 // Lunghezza del mondo  
	WORLD_DEPTH      = 100 // Profondità del mondo
)

// Tipi di materiali
const (
	AIR = iota
	DIRT
	STONE
	COAL
	IRON
	GOLD
	DIAMOND
	BEDROCK = 10
)

// Colori ANSI per output colorato
const (
	RESET  = "\033[0m"
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"
	PURPLE = "\033[35m"
	CYAN   = "\033[36m"
	WHITE  = "\033[37m"
	GRAY   = "\033[90m"
	BOLD   = "\033[1m"
)

// Strutture dati
type Vector2D struct {
	X, Y float64
}

type BlockPosition struct {
	Material int
	X, Y, Z  int
}

type PerlinGenerator struct {
	gridVectors    [GRID_SLOTS + 1][GRID_SLOTS + 1]Vector2D
	perlinMaps     [4][GRID_SLOTS][GRID_SLOTS][PIXELS_PER_SLOT][PIXELS_PER_SLOT]float64
	finalMap       [GRID_SLOTS][GRID_SLOTS][PIXELS_PER_SLOT][PIXELS_PER_SLOT]float64
	worldHeightMap [GRID_SLOTS * PIXELS_PER_SLOT][GRID_SLOTS * PIXELS_PER_SLOT]float64
}

type WorldGenerator struct {
	world      [WORLD_HEIGHT][WORLD_LENGTH][WORLD_DEPTH]int
	veins      map[BlockPosition]int
	perlinGen  *PerlinGenerator
}

// Configurazione minerali
type MineralConfig struct {
	probability int
	veinSize    int
	minDepth    int
	maxDepth    int
}

var mineralConfigs = map[int]MineralConfig{
	COAL:    {probability: 13, veinSize: 20, minDepth: 2, maxDepth: 40},
	IRON:    {probability: 8, veinSize: 15, minDepth: 2, maxDepth: 40},
	GOLD:    {probability: 5, veinSize: 10, minDepth: 2, maxDepth: 30},
	DIAMOND: {probability: 3, veinSize: 5, minDepth: 2, maxDepth: 15},
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	printHeader()
	
	// Inizializza generatori
	perlinGen := NewPerlinGenerator()
	worldGen := NewWorldGenerator(perlinGen)
	
	// Genera mappa Perlin
	fmt.Print(CYAN + "🌊 Generando mappa Perlin..." + RESET)
	perlinGen.Generate()
	fmt.Println(GREEN + " ✓ Completato!" + RESET)
	
	// Menu principale
	showMainMenu(perlinGen, worldGen)
}

func printHeader() {
	fmt.Println(BOLD + BLUE + "╔══════════════════════════════════════════════════════════╗" + RESET)
	fmt.Println(BOLD + BLUE + "║" + RESET + BOLD + "            🌍 GENERATORE MONDO PERLIN 🌍              " + BOLD + BLUE + "║" + RESET)
	fmt.Println(BOLD + BLUE + "╚══════════════════════════════════════════════════════════╝" + RESET)
	fmt.Println()
}

func showMainMenu(perlinGen *PerlinGenerator, worldGen *WorldGenerator) {
	for {
		fmt.Println(BOLD + "\n📋 MENU PRINCIPALE:" + RESET)
		fmt.Println("1️⃣  " + CYAN + "Visualizza mappa Perlin e grafici altezza" + RESET)
		fmt.Println("2️⃣  " + YELLOW + "Genera mondo e mostra statistiche minerali" + RESET)
		fmt.Println("3️⃣  " + RED + "Esci" + RESET)
		fmt.Print("\n" + BOLD + "Scelta: " + RESET)
		
		var choice int
		fmt.Scan(&choice)
		
		switch choice {
		case 1:
			showPerlinMenu(perlinGen)
		case 2:
			generateWorldMenu(worldGen)
		case 3:
			fmt.Println(GREEN + "\n👋 Arrivederci!" + RESET)
			os.Exit(0)
		default:
			fmt.Println(RED + "❌ Scelta non valida!" + RESET)
		}
	}
}

func showPerlinMenu(perlinGen *PerlinGenerator) {
	fmt.Println(CYAN + "\n🗺️  VISUALIZZAZIONE MAPPA PERLIN" + RESET)
	perlinGen.DrawMap()
	
	fmt.Println(BOLD + "\n📊 GRAFICI ALTEZZA:" + RESET)
	fmt.Println("1️⃣  Grafici per righe")
	fmt.Println("2️⃣  Grafici per colonne")
	fmt.Print("\nScelta: ")
	
	var choice int
	fmt.Scan(&choice)
	
	if choice == 1 || choice == 2 {
		perlinGen.DrawHeightGraphics(choice)
	} else {
		fmt.Println(RED + "❌ Scelta non valida!" + RESET)
	}
}

func generateWorldMenu(worldGen *WorldGenerator) {
	fmt.Print(YELLOW + "⛏️  Generando mondo..." + RESET)
	worldGen.GenerateWorld()
	fmt.Println(GREEN + " ✓ Completato!" + RESET)
	
	worldGen.ShowStatistics()
	showWorldVisualizationMenu(worldGen)
}

func showWorldVisualizationMenu(worldGen *WorldGenerator) {
	fmt.Println(BOLD + "\n🔍 VISUALIZZAZIONE MONDO:" + RESET)
	fmt.Print("Minerale da visualizzare (coal/iron/gold/diamond): ")
	
	var mineralName string
	fmt.Scan(&mineralName)
	
	mineralType := getMineralType(mineralName)
	if mineralType == -1 {
		fmt.Println(RED + "❌ Minerale non valido!" + RESET)
		return
	}
	
	for {
		fmt.Println(BOLD + "\n📐 TIPO DI SEZIONE:" + RESET)
		fmt.Println("1️⃣  Sezione orizzontale (vista dall'alto)")
		fmt.Println("2️⃣  Sezione verticale frontale")
		fmt.Println("3️⃣  Sezione verticale laterale")
		fmt.Println("4️⃣  Torna al menu principale")
		fmt.Print("\nScelta: ")
		
		var choice int
		fmt.Scan(&choice)
		
		if choice >= 1 && choice <= 3 {
			worldGen.PrintSection(mineralType, choice)
		} else if choice == 4 {
			break
		} else {
			fmt.Println(RED + "❌ Scelta non valida!" + RESET)
		}
	}
}

// Costruttori
func NewPerlinGenerator() *PerlinGenerator {
	pg := &PerlinGenerator{}
	return pg
}

func NewWorldGenerator(perlinGen *PerlinGenerator) *WorldGenerator {
	return &WorldGenerator{
		perlinGen: perlinGen,
		veins:     make(map[BlockPosition]int),
	}
}

// Metodi PerlinGenerator
func (pg *PerlinGenerator) Generate() {
	pg.createGridVectors()
	pg.createPerlinMaps()
	pg.interpolateMaps()
	pg.resizeValues()
	pg.createWorldHeightMap()
}

func (pg *PerlinGenerator) createGridVectors() {
	for i := 0; i <= GRID_SLOTS; i++ {
		for j := 0; j <= GRID_SLOTS; j++ {
			// Genera vettore casuale normalizzato
			angle := rand.Float64() * 2 * math.Pi
			pg.gridVectors[i][j] = Vector2D{
				X: math.Cos(angle),
				Y: math.Sin(angle),
			}
		}
	}
}

func (pg *PerlinGenerator) createPerlinMaps() {
	unitVec := 1.0 / float64(PIXELS_PER_SLOT)
	
	for i := 0; i < GRID_SLOTS; i++ {
		for j := 0; j < GRID_SLOTS; j++ {
			for k := 0; k < PIXELS_PER_SLOT; k++ {
				for v := 0; v < PIXELS_PER_SLOT; v++ {
					// Calcola offset per ogni corner
					offsetY := float64(k) * unitVec
					offsetX := float64(v) * unitVec
					
					// Corner vectors
					corners := []struct {
						gridI, gridJ int
						offsetX, offsetY float64
						mapIndex int
					}{
						{i, j, -offsetX, offsetY, 0},                    // up-left
						{i, j+1, 1-offsetX, offsetY, 1},                // up-right  
						{i+1, j, -offsetX, -(1-offsetY), 2},            // down-left
						{i+1, j+1, 1-offsetX, -(1-offsetY), 3},        // down-right
					}
					
					for _, corner := range corners {
						pg.perlinMaps[corner.mapIndex][i][j][k][v] = pg.dotProduct(
							pg.gridVectors[corner.gridI][corner.gridJ],
							Vector2D{corner.offsetX, corner.offsetY},
						)
					}
				}
			}
		}
	}
}

func (pg *PerlinGenerator) dotProduct(v1, v2 Vector2D) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func (pg *PerlinGenerator) interpolateMaps() {
	for i := 0; i < GRID_SLOTS; i++ {
		for j := 0; j < GRID_SLOTS; j++ {
			for k := 0; k < PIXELS_PER_SLOT; k++ {
				for v := 0; v < PIXELS_PER_SLOT; v++ {
					u := float64(k) / float64(PIXELS_PER_SLOT-1)
					w := float64(v) / float64(PIXELS_PER_SLOT-1)
					
					pg.finalMap[i][j][k][v] = pg.bilinearInterpolation(
						pg.perlinMaps[0][i][j][k][v], // up-left
						pg.perlinMaps[1][i][j][k][v], // up-right
						pg.perlinMaps[2][i][j][k][v], // down-left
						pg.perlinMaps[3][i][j][k][v], // down-right
						u, w,
					)
				}
			}
		}
	}
}

func (pg *PerlinGenerator) bilinearInterpolation(x1, x2, x3, x4, u, w float64) float64 {
	return x1*(1-u)*(1-w) + x2*u*(1-w) + x3*(1-u)*w + x4*u*w
}

func (pg *PerlinGenerator) resizeValues() {
	for i := 0; i < GRID_SLOTS; i++ {
		for j := 0; j < GRID_SLOTS; j++ {
			for k := 0; k < PIXELS_PER_SLOT; k++ {
				for v := 0; v < PIXELS_PER_SLOT; v++ {
					val := pg.finalMap[i][j][k][v] * 10.0
					if val > 9 {
						val = 9
					} else if val < -9 {
						val = -9
					}
					pg.finalMap[i][j][k][v] = val
				}
			}
		}
	}
}

func (pg *PerlinGenerator) createWorldHeightMap() {
	for i := 0; i < GRID_SLOTS; i++ {
		for j := 0; j < GRID_SLOTS; j++ {
			for k := 0; k < PIXELS_PER_SLOT; k++ {
				for v := 0; v < PIXELS_PER_SLOT; v++ {
					destI := i*PIXELS_PER_SLOT + k
					destJ := j*PIXELS_PER_SLOT + v
					pg.worldHeightMap[destI][destJ] = pg.finalMap[i][j][k][v] + 9
				}
			}
		}
	}
}

func (pg *PerlinGenerator) DrawMap() {
	fmt.Println(BOLD + "\n🗺️  MAPPA PERLIN NOISE:" + RESET)
	fmt.Println(strings.Repeat("─", 80))
	
	for i := 0; i < GRID_SLOTS; i++ {
		for k := 0; k < PIXELS_PER_SLOT; k++ {
			for j := 0; j < GRID_SLOTS; j++ {
				for v := 0; v < PIXELS_PER_SLOT; v++ {
					val := pg.finalMap[i][j][k][v]
					color := pg.getHeightColor(val)
					if val < 0 {
						fmt.Printf("%s%3.0f%s ", color, val, RESET)
					} else {
						fmt.Printf("%s+%2.0f%s ", color, val, RESET)
					}
				}
			}
			fmt.Println()
		}
	}
}

func (pg *PerlinGenerator) getHeightColor(val float64) string {
	switch {
	case val <= -6:
		return BLUE    // Molto basso (acqua profonda)
	case val <= -3:
		return CYAN    // Basso (acqua)
	case val <= 0:
		return GREEN   // Livello mare
	case val <= 3:
		return YELLOW  // Collina
	case val <= 6:
		return RED     // Montagna
	default:
		return WHITE   // Picco
	}
}

func (pg *PerlinGenerator) DrawHeightGraphics(mode int) {
	fmt.Printf(BOLD + "\n📊 GRAFICO ALTEZZE (Modalità %d):\n" + RESET, mode)
	fmt.Println(strings.Repeat("─", 60))
	
	if mode == 1 {
		// Grafici per righe
		for i := 0; i < GRID_SLOTS; i++ {
			for k := 0; k < PIXELS_PER_SLOT; k++ {
				fmt.Printf(GRAY + "Riga %d-%d: " + RESET, i, k)
				for j := 0; j < GRID_SLOTS; j++ {
					for v := 0; v < PIXELS_PER_SLOT; v++ {
						height := int(pg.finalMap[i][j][k][v])
						color := pg.getHeightColor(float64(height))
						barLength := int(math.Abs(float64(height)))
						if barLength > 0 {
							fmt.Print(color + strings.Repeat("█", barLength) + RESET)
						}
					}
				}
				fmt.Println()
			}
		}
	} else {
		// Grafici per colonne
		for j := 0; j < GRID_SLOTS; j++ {
			for v := 0; v < PIXELS_PER_SLOT; v++ {
				fmt.Printf(GRAY + "Col %d-%d: " + RESET, j, v)
				for i := 0; i < GRID_SLOTS; i++ {
					for k := 0; k < PIXELS_PER_SLOT; k++ {
						height := int(pg.finalMap[i][j][k][v])
						color := pg.getHeightColor(float64(height))
						barLength := int(math.Abs(float64(height)))
						if barLength > 0 {
							fmt.Print(color + strings.Repeat("█", barLength) + RESET)
						}
					}
				}
				fmt.Println()
			}
		}
	}
}

// Metodi WorldGenerator
func (wg *WorldGenerator) GenerateWorld() {
	wg.veins = make(map[BlockPosition]int)
	
	for i := 0; i < WORLD_HEIGHT; i++ {
		for j := 0; j < WORLD_LENGTH; j++ {
			for k := 0; k < WORLD_DEPTH; k++ {
				wg.world[i][j][k] = wg.generateBlock(i, j, k)
			}
		}
	}
}

func (wg *WorldGenerator) generateBlock(height, length, depth int) int {
	// Bedrock al livello 0
	if height == 0 {
		return BEDROCK
	}
	
	// Zona superficie con Perlin noise
	if height > 40 {
		if length < GRID_SLOTS*PIXELS_PER_SLOT && depth < GRID_SLOTS*PIXELS_PER_SLOT {
			surfaceHeight := int(wg.perlinGen.worldHeightMap[length][depth])
			if (height - 40) <= surfaceHeight {
				return DIRT
			}
		}
		return AIR
	}
	
	// Generazione sotterranea
	return wg.generateUndergroundBlock(height, length, depth)
}

func (wg *WorldGenerator) generateUndergroundBlock(height, length, depth int) int {
	// Controlla se c'è un minerale vicino
	if nearbyMineral, found := wg.findNearbyMineral(height, length, depth); found {
		return wg.generateMineralOrStone(nearbyMineral, height, length, depth)
	}
	
	// Generazione casuale di minerali
	randomValue := rand.Intn(600)
	
	for mineral, config := range mineralConfigs {
		if height >= config.minDepth && height <= config.maxDepth && randomValue < config.probability {
			wg.createNewVein(mineral, depth, length, height)
			return mineral
		}
	}
	
	return STONE
}

func (wg *WorldGenerator) findNearbyMineral(height, length, depth int) (BlockPosition, bool) {
	directions := []struct{ dh, dl, dd int }{
		{-1, 0, 0}, {0, -1, 0}, {0, 0, -1},
	}
	
	for _, dir := range directions {
		nh, nl, nd := height+dir.dh, length+dir.dl, depth+dir.dd
		if nh >= 0 && nl >= 0 && nd >= 0 &&
		   nh < WORLD_HEIGHT && nl < WORLD_LENGTH && nd < WORLD_DEPTH {
			material := wg.world[nh][nl][nd]
			if wg.isMaterial(material) {
				return BlockPosition{material, nd, nl, nh}, true
			}
		}
	}
	return BlockPosition{}, false
}

func (wg *WorldGenerator) isMaterial(material int) bool {
	return material == COAL || material == IRON || material == GOLD || material == DIAMOND
}

func (wg *WorldGenerator) generateMineralOrStone(nearbyBlock BlockPosition, height, length, depth int) int {
	if rand.Intn(3) == 0 {
		return STONE
	}
	
	veinKey, inVein := wg.findVein(nearbyBlock)
	if inVein {
		config := mineralConfigs[nearbyBlock.Material]
		if wg.veins[veinKey] < config.veinSize {
			newBlock := BlockPosition{nearbyBlock.Material, depth, length, height}
			wg.updateVein(newBlock)
			return nearbyBlock.Material
		}
	}
	
	return STONE
}

func (wg *WorldGenerator) createNewVein(material, x, y, z int) {
	vein := BlockPosition{material, x, y, z}
	wg.veins[vein] = 1
}

func (wg *WorldGenerator) findVein(block BlockPosition) (BlockPosition, bool) {
	for veinKey := range wg.veins {
		if veinKey.Material == block.Material {
			if wg.isNearby(veinKey, block, 5) {
				return veinKey, true
			}
		}
	}
	return BlockPosition{}, false
}

func (wg *WorldGenerator) isNearby(pos1, pos2 BlockPosition, limit int) bool {
	dx := int(math.Abs(float64(pos1.X - pos2.X)))
	dy := int(math.Abs(float64(pos1.Y - pos2.Y)))
	dz := int(math.Abs(float64(pos1.Z - pos2.Z)))
	return dx <= limit && dy <= limit && dz <= limit
}

func (wg *WorldGenerator) updateVein(block BlockPosition) {
	veinKey, found := wg.findVein(block)
	if found {
		oldValue := wg.veins[veinKey]
		wg.veins[block] = oldValue + 1
		delete(wg.veins, veinKey)
	}
}

func (wg *WorldGenerator) ShowStatistics() {
	counts := make(map[int]int)
	
	for i := 0; i < WORLD_HEIGHT; i++ {
		for j := 0; j < WORLD_LENGTH; j++ {
			for k := 0; k < WORLD_DEPTH; k++ {
				counts[wg.world[i][j][k]]++
			}
		}
	}
	
	fmt.Println(BOLD + "\n📊 STATISTICHE MONDO:" + RESET)
	fmt.Println(strings.Repeat("─", 40))
	
	materials := []struct {
		material int
		name     string
		color    string
		icon     string
	}{
		{COAL, "Carbone", GRAY, "⚫"},
		{IRON, "Ferro", RED, "🔴"},
		{GOLD, "Oro", YELLOW, "🟡"},
		{DIAMOND, "Diamante", CYAN, "💎"},
	}
	
	for _, mat := range materials {
		count := counts[mat.material]
		percentage := float64(count) / float64(WORLD_HEIGHT*WORLD_LENGTH*WORLD_DEPTH) * 100
		fmt.Printf("%s %s%s%s: %s%d%s blocchi (%.2f%%)\n",
			mat.icon, mat.color, mat.name, RESET, BOLD, count, RESET, percentage)
	}
	
	fmt.Printf("\n🏔️  %sAltri materiali%s: %s%d%s blocchi\n",
		GRAY, RESET, BOLD, 
		counts[STONE]+counts[DIRT]+counts[BEDROCK]+counts[AIR], RESET)
}

func (wg *WorldGenerator) PrintSection(mineralType, sectionType int) {
	fmt.Printf(BOLD + "\n🔍 SEZIONE TIPO %d - %s:\n" + RESET, 
		sectionType, getMaterialName(mineralType))
	
	legend := map[int]string{
		AIR:     " ",
		BEDROCK: GRAY + "@" + RESET,
		STONE:   GRAY + "▓" + RESET,
		DIRT:    YELLOW + "░" + RESET,
		COAL:    GRAY + "●" + RESET,
		IRON:    RED + "●" + RESET,
		GOLD:    YELLOW + "●" + RESET,
		DIAMOND: CYAN + "♦" + RESET,
	}
	
	// Stampa legenda
	fmt.Println(BOLD + "Legenda:" + RESET)
	fmt.Printf("  Aria: [%s]  Bedrock: [%s]  Pietra: [%s]  Terra: [%s]\n",
		legend[AIR], legend[BEDROCK], legend[STONE], legend[DIRT])
	fmt.Printf("  %s: [%s]  (evidenziato: %s[#]%s)\n\n",
		getMaterialName(mineralType), legend[mineralType], 
		getMaterialColor(mineralType), RESET)
	
	switch sectionType {
	case 1: // Sezione orizzontale
		wg.printHorizontalSections(mineralType, legend)
	case 2: // Sezione verticale frontale
		wg.printVerticalFrontSections(mineralType, legend)
	case 3: // Sezione verticale laterale
		wg.printVerticalSideSections(mineralType, legend)
	}
}

func (wg *WorldGenerator) printHorizontalSections(mineralType int, legend map[int]string) {
	for i := WORLD_HEIGHT - 1; i >= 0; i-- {
		fmt.Printf(BOLD + "Livello %d:\n" + RESET, i)
		for j := 0; j < WORLD_LENGTH; j++ {
			for k := 0; k < WORLD_DEPTH; k++ {
				material := wg.world[i][j][k]
				if material == mineralType {
					fmt.Print(getMaterialColor(mineralType) + "#" + RESET)
				} else {
					fmt.Print(legend[material])
				}
			}
			fmt.Println()
		}
		fmt.Println(strings.Repeat("─", 60))
	}
}

func (wg *WorldGenerator) printVerticalFrontSections(mineralType int, legend map[int]string) {
	for k := 0; k < WORLD_DEPTH; k++ {
		fmt.Printf(BOLD + "Sezione frontale %d:\n" + RESET, k)
		for i := WORLD_HEIGHT - 1; i >= 0; i-- {
			for j := 0; j < WORLD_LENGTH; j++ {
				material := wg.world[i][j][k]
				if material == mineralType {
					fmt.Print(getMaterialColor(mineralType) + "#" + RESET)
				} else {
					fmt.Print(legend[material])
				}
			}
			fmt.Println()
		}
		fmt.Println(strings.Repeat("─", 60))
	}
}

func (wg *WorldGenerator) printVerticalSideSections(mineralType int, legend map[int]string) {
	for j := 0; j < WORLD_LENGTH; j++ {
		fmt.Printf(BOLD + "Sezione laterale %d:\n" + RESET, j)
		for i := WORLD_HEIGHT - 1; i >= 0; i-- {
			for k := 0; k < WORLD_DEPTH; k++ {
				material := wg.world[i][j][k]
				if material == mineralType {
					fmt.Print(getMaterialColor(mineralType) + "#" + RESET)
				} else {
					fmt.Print(legend[material])
				}
			}
			fmt.Println()
		}
		fmt.Println(strings.Repeat("─", 60))
	}
}

// Funzioni helper
func getMineralType(name string) int {
	minerals := map[string]int{
		"coal":    COAL,
		"iron":    IRON,
		"gold":    GOLD,
		"diamond": DIAMOND,
	}
	if material, exists := minerals[name]; exists {
		return material
	}
	return -1
}

func getMaterialName(material int) string {
	names := map[int]string{
		COAL:    "Carbone",
		IRON:    "Ferro", 
		GOLD:    "Oro",
		DIAMOND: "Diamante",
	}
	return names[material]
}

func getMaterialColor(material int) string {
	colors := map[int]string{
		COAL:    GRAY,
		IRON:    RED,
		GOLD:    YELLOW,
		DIAMOND: CYAN,
	}
	return colors[material]
}