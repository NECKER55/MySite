package main

import (
	"fmt"
	"math"
	"math/rand"
)

const n_slot int = 10      // numero slot della griglia (in totale saranno: n_slot*n_slot)
const n_pixel_in_slot = 10 // misura lato non l'area (in ogni slot ci saranno: n_pixel_in_slot * n_pixel_in_slot)

var total_vectors [n_slot + 1][n_slot + 1][2]float64 // vettori creati in ogni vertice della griglia

var perlin_map_up_l [n_slot][n_slot][n_pixel_in_slot][n_pixel_in_slot]float64   //mappa di perlin ottenuta con il prodotto scalare tra ogni pixel e il vettore in alto a sinistra della sua griglia
var perlin_map_up_r [n_slot][n_slot][n_pixel_in_slot][n_pixel_in_slot]float64   //mappa di perlin ottenuta con il prodotto scalare tra ogni pixel e il vettore in alto a destra della sua griglia
var perlin_map_down_l [n_slot][n_slot][n_pixel_in_slot][n_pixel_in_slot]float64 //mappa di perlin ottenuta con il prodotto scalare tra ogni pixel e il vettore in basso a sinistra della sua griglia
var perlin_map_down_r [n_slot][n_slot][n_pixel_in_slot][n_pixel_in_slot]float64 //mappa di perlin ottenuta con il prodotto scalare tra ogni pixel e il vettore in basso a destra della sua griglia

var perlin_map [n_slot][n_slot][n_pixel_in_slot][n_pixel_in_slot]float64 // mappa di perlin finale
var perlin_map_for_world [n_slot * n_pixel_in_slot][n_slot * n_pixel_in_slot]float64

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
const lenght_blocks int = 100 // numero di blocchi per lunghezza
const deep_blocks int = 100   // numero di blocchi in profondità
const h_blocks int = 70       // numero di blocchi in altezza per raggiungere lo zero nella mappa di perlin

const air = 0 // legenda dei matriali
const bedrock = 10
const dirt int = 1
const stone int = 2
const coal int = 3
const iron int = 4
const gold int = 5
const diamond int = 6

type block_position = struct {
	material int
	pos_x    int
	pos_y    int
	pos_z    int
}

var total_veins map[block_position]int // memoria delle vene minerarie totali
var current_vein int

var world [h_blocks][lenght_blocks][deep_blocks]int // array che genera il mondo secondo la logia seguente: per ogni livello in altezza si sviluppa un piano dove per ogni blocco in lunghezza si sviluppa una fila di blocchi, per comodità immaginiamo il mondo come una pila di piani i quali sono formati da colonne (lenght_blocks) e righe (deep_blocks)

func main() {

	create_total_vector()
	create_perlin_maps()
	perlin_map_interpolation()
	resizing_values()
	create_perlin_map_for_world()

	var a int
	fmt.Println("- premere 1 per visualizzare la mappa di perlin e il grafico dei dislivelli")
	fmt.Println("- premere 2 per visualizzare il numero di minerali e ogni strato del mondo")
	fmt.Scan(&a)

	if a == 1 {
		draw_perlin_map()
		graphic_of_hight()
	} else if a == 2 {
		assignment_of_blocks()
		data()
		print_levels()
	}
}

func create_total_vector() { //crea i vettori in ogni punto della griglia principale
	for i := 0; i <= n_slot; i++ {
		for j := 0; j <= n_slot; j++ {

			selection := rand.Intn(4)
			random_pos1 := rand.Float64()
			random_pos2 := rand.Float64()
			random_neg1 := rand.Float64() * (-1)
			random_neg2 := rand.Float64() * (-1) // servono ad ottenere qualsiasi tipo di inclinazione del vettore

			switch {
			case selection == 0:
				total_vectors[i][j][0] = random_pos1
				total_vectors[i][j][1] = random_pos2
			case selection == 1:
				total_vectors[i][j][0] = random_pos1
				total_vectors[i][j][1] = random_neg1
			case selection == 2:
				total_vectors[i][j][0] = random_neg1
				total_vectors[i][j][1] = random_pos1
			case selection == 3:
				total_vectors[i][j][0] = random_neg1
				total_vectors[i][j][1] = random_neg2
			}
		}
	}
}

func create_perlin_maps() { // crea le 4 mappe di perlin

	var unit_vec float64 = 1 / n_pixel_in_slot

	for i := 0; i < n_slot; i++ {
		for j := 0; j < n_slot; j++ {
			for k := 0; k < n_pixel_in_slot; k++ {
				for v := 0; v < n_pixel_in_slot; v++ {
					var comp_y_up_l float64 = float64(k) * unit_vec
					var comp_x_up_l float64 = -(float64(v) * unit_vec)

					var comp_y_up_r float64 = float64(k) * unit_vec
					var comp_x_up_r float64 = 1 - (float64(v) * unit_vec)

					var comp_y_down_l float64 = -(1 - (float64(k) * unit_vec))
					var comp_x_down_l float64 = -(float64(v) * unit_vec)

					var comp_y_down_r float64 = -(1 - (float64(k) * unit_vec))
					var comp_x_down_r float64 = 1 - (float64(v) * unit_vec)

					perlin_map_up_l[i][j][k][v] = dot_product(total_vectors[i][j][0], total_vectors[i][j][1], comp_x_up_l, comp_y_up_l)
					perlin_map_up_r[i][j][k][v] = dot_product(total_vectors[i][j+1][0], total_vectors[i][j+1][1], comp_x_up_r, comp_y_up_r)
					perlin_map_down_l[i][j][k][v] = dot_product(total_vectors[i+1][j][0], total_vectors[i+1][j][1], comp_x_down_l, comp_y_down_l)
					perlin_map_down_r[i][j][k][v] = dot_product(total_vectors[i+1][j+1][0], total_vectors[i+1][j+1][1], comp_x_down_r, comp_y_down_r)
				}
			}
		}
	}
}

func dot_product(vec_x, vec_y, comp_x, comp_y float64) float64 { //fa il prodotto scalare tra il vettore principale e i vettori dei pixel

	var dot_product float64

	dot_product = (vec_x * comp_x) + (vec_y * comp_y)
	return dot_product
}

//func perlin_map_interpolation() {
//	for i := 0; i < n_slot; i++ {
//		for j := 0; j < n_slot; j++ {
//			for k := 0; k < n_pixel_in_slot; k++ {
//				for v := 0; v < n_pixel_in_slot; v++ {
//					perlin_map[i][j][k][v] = (perlin_map_up_l[i][j][k][v] + perlin_map_up_r[i][j][k][v] + perlin_map_down_l[i][j][k][v] + perlin_map_down_r[i][j][k][v]) / 4
//				}
//			}
//		}
//	}
//fmt.Print(perlin_map)
//}

func perlin_map_interpolation() { // unisce le 4 mappe in una
	for i := 0; i < n_slot; i++ {
		for j := 0; j < n_slot; j++ {
			for k := 0; k < n_pixel_in_slot; k++ {
				for v := 0; v < n_pixel_in_slot; v++ {
					// Utilizza l'interpolazione bilineare per ottenere una transizione più fluida
					u := float64(k) / float64(n_pixel_in_slot-1)
					w := float64(v) / float64(n_pixel_in_slot-1)

					perlin_map[i][j][k][v] = bilinearInterpolation(
						perlin_map_up_l[i][j][k][v],
						perlin_map_up_r[i][j][k][v],
						perlin_map_down_l[i][j][k][v],
						perlin_map_down_r[i][j][k][v],
						u, w,
					)
				}
			}
		}
	}
}

func perlin_map_interpolation1() { // unisce le 4 mappe in una
	for i := 0; i < n_slot; i++ {
		for j := 0; j < n_slot; j++ {
			for k := 0; k < n_pixel_in_slot; k++ {
				for v := 0; v < n_pixel_in_slot; v++ {
					// Utilizza l'interpolazione bilineare per ottenere una transizione più fluida
					u := float64(k) / float64(n_pixel_in_slot-1)
					w := float64(v) / float64(n_pixel_in_slot-1)

					//					perlin_map[i][j][k][v] = bilinearInterpolation(
					//						perlin_map_up_l[i][j][k][v],
					//						perlin_map_up_r[i][j][k][v],
					//						perlin_map_down_l[i][j][k][v],
					//						perlin_map_down_r[i][j][k][v],
					//						u, w,
					//					)

					perlin_map[i][j][k][v] = cosineInterpolation(
						cosineInterpolation(
							perlin_map_up_l[i][j][k][v],
							perlin_map_up_r[i][j][k][v],
							u,
						),
						cosineInterpolation(
							perlin_map_down_l[i][j][k][v],
							perlin_map_down_r[i][j][k][v],
							u,
						),
						w,
					)
				}
			}
		}
	}
	resizing_values()
	draw_perlin_map()
	perlin_map_interpolation()
}

func bilinearInterpolation(x1, x2, x3, x4, u, w float64) float64 { // smussa le differenze tra ogni mappa
	p1 := x1 * (1 - u) * (1 - w)
	p2 := x2 * u * (1 - w)
	p3 := x3 * (1 - u) * w
	p4 := x4 * u * w
	return p1 + p2 + p3 + p4
}

func cosineInterpolation(y0, y1, mu float64) float64 {
	mu2 := (1 - math.Cos(mu*math.Pi)) / 2
	return y0*(1-mu2) + y1*mu2
}

func resizing_values() { // assegna valori ad ogni pixel per rendere la mappa leggibile

	for i := 0; i < n_slot; i++ {
		for j := 0; j < n_slot; j++ {
			for k := 0; k < n_pixel_in_slot; k++ {
				for v := 0; v < n_pixel_in_slot; v++ {

					perlin_map[i][j][k][v] = float64(int(perlin_map[i][j][k][v] * 10.0))

					if perlin_map[i][j][k][v] > 9 {
						perlin_map[i][j][k][v] = 9
					}
					if perlin_map[i][j][k][v] < -9 {
						perlin_map[i][j][k][v] = -9
					}
				}
			}
		}
	}
}

func create_perlin_map_for_world() { // inserisce la mappa in un'array ordinato e non diviso in slot
	for i := 0; i < n_slot; i++ {
		for j := 0; j < n_slot; j++ {
			for k := 0; k < n_pixel_in_slot; k++ {
				for v := 0; v < n_pixel_in_slot; v++ { // Calcolo delle posizioni negli array di destinazione
					dest_i := i*n_pixel_in_slot + k
					dest_j := j*n_pixel_in_slot + v

					perlin_map_for_world[dest_i][dest_j] = perlin_map[i][j][k][v] // Copia del valore dall'array di origine a quello di destinazione
					perlin_map_for_world[dest_i][dest_j] += 9

				}
			}
		}
	}

}

func draw_perlin_map() { // stampa la mappa a schermo

	for i := 0; i < n_slot; i++ {
		for k := 0; k < n_pixel_in_slot; k++ {
			for j := 0; j < n_slot; j++ {
				for v := 0; v < n_pixel_in_slot; v++ {
					if int(int(perlin_map[i][j][k][v])) < 0 {
						fmt.Print(perlin_map[i][j][k][v], "|")
					} else {
						fmt.Print("+", perlin_map[i][j][k][v], "|")
					}
				}
			}
			fmt.Println()
		}
	}
	fmt.Println()
	fmt.Println()
	fmt.Println()
}

func graphic_of_hight() {

	var a int

	fmt.Println("premi:")
	fmt.Println("- 1 per vedere i dislivelli di ogni riga")
	fmt.Println("- 1 per vedere i dislivelli di ogni colonna")
	fmt.Scan(&a)

	// stampa il dislivello per ogni riga

	if a == 1 {
		for i := 0; i < n_slot; i++ {
			for k := 0; k < n_pixel_in_slot; k++ {
				for j := 0; j < n_slot; j++ {
					for v := 0; v < n_pixel_in_slot; v++ {
						for h := -9; h < int(perlin_map[i][j][k][v]); h++ {
							fmt.Print("#")
						}
						fmt.Println()
					}
				}
				fmt.Println()
				fmt.Println()
			}
		}

		fmt.Println()
		fmt.Println("------------------------------------------------------------------------------------------------")
		fmt.Println()
	} else if a == 2 { // stampa il dislivello per ogni colonna
		for j := 0; j < n_slot; j++ {
			for v := 0; v < n_pixel_in_slot; v++ {
				for i := 0; i < n_slot; i++ {
					for k := 0; k < n_pixel_in_slot; k++ {
						for h := -9; h < int(perlin_map[i][j][k][v]); h++ {
							fmt.Print("#")
						}
						fmt.Println()
					}
				}
				fmt.Println()
				fmt.Println()
			}
		}
	} else {
		fmt.Println("ERROR")
	}

}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func assignment_of_blocks() { // scorre l'array che contiene il mondo e a ogni spazio alloca un blocco

	total_veins = make(map[block_position]int) // alloca in memoria spazio per le varie vene di minerali
	//fmt.Print(total_veins)

	for i := 0; i < h_blocks; i++ {
		for j := 0; j < lenght_blocks; j++ {
			for k := 0; k < deep_blocks; k++ {
				world[i][j][k] = casual_generation_world(i, j, k)
			}
		}

	}
}

func casual_generation_world(high, lenght, column int) int { // calcola quale blocco deve allocare

	var current_block block_position // struct che contiene materiale e posizione del blocco che dovrà essere generato
	current_block.pos_x = column
	current_block.pos_y = lenght
	current_block.pos_z = high

	var coal_probability int = 13   // percentuale di probabilità di generare del carbone
	var iron_probability int = 8    // percentuale di probabilità di generare del ferro
	var gold_probability int = 5    // percentuale di probabilità di generare dell'oro
	var diamond_probability int = 3 // percentuale di probabilità di generare del diamante

	var coal_vein int = 20   // numero massimo di minerale per vena
	var iron_vein int = 15   // numero massimo di minerale per vena
	var gold_vein int = 10   // numero massimo di minerale per vena
	var diamond_vein int = 5 // numero massimo di minerale per vena

	stone_level := [2]int{1, 40}   // zona di generazione della pietra
	coal_level := [2]int{2, 40}    // zona di generazione del carbone
	iron_level := [2]int{2, 40}    // zona di generazione del ferro
	gold_level := [2]int{2, 30}    // zona di generazione dell'oro
	diamond_level := [2]int{2, 15} // zona di generazione del diamante

	_, _, position_y_previous_block, position_x_previous_block, previous_block_exist, previous_block_is_mineral := mineral_recognition(high, lenght, column) // assegna la posizione x e y del blocco precedente che ci servirà per calcolare la probabilità di generazione del minerale, se previous_block_exist è vero allora il blocco da generare non è il primo e se previous_block_is_mineral è vero allora il blocco precedente è un minerale

	if high == 0 { // genera al livello zero uno strato di un materiale ch non può essere distrutto
		current_block.material = bedrock
	} else if high > stone_level[1] { // se l'altezza del mondo supera l'altezza massima della generazione di pietra genera solo terra
		if lenght <= n_slot*n_pixel_in_slot && column <= n_slot*n_pixel_in_slot { // permette la generazione del suolo in base alla mappa di perlin
			if (high - stone_level[1]) <= int(perlin_map_for_world[lenght][column]) {
				current_block.material = dirt
			} else {
				current_block.material = air
			}
		} else {
			current_block.material = air
		}
	} else if previous_block_exist && previous_block_is_mineral == false && veins_nearby(high, lenght, column) == false { // verifica se i blocchi vicini alla generazione siano pietra o bedrock o terra e se il blocco da generare non è il primo
		random := rand.Intn(600)

		if high >= diamond_level[0] && high <= diamond_level[1] && diamond_probability >= random { // verifica le condizioni in cui può essere generato il minerale
			current_block.material = diamond

		} else if high >= gold_level[0] && high <= gold_level[1] && gold_probability >= random { // verifica le condizioni in cui può essere generato il minerale
			current_block.material = gold

		} else if high >= iron_level[0] && high <= iron_level[1] && iron_probability >= random { // verifica le condizioni in cui può essere generato il minerale
			current_block.material = iron

		} else if high >= coal_level[0] && high <= coal_level[1] && coal_probability >= random { // verifica le condizioni in cui può essere generato il minerale
			current_block.material = coal

		} else { // se le precedenti condizioni non sono verificate, genera della pietra
			current_block.material = stone
		}
		if current_block.material != stone && current_block.material != bedrock && current_block.material != dirt { // se viene generato un minerale, si crea una vena mineraria nuova
			create_new_vein(current_block.material, column, lenght, high)
		}

	} else if previous_block_exist && previous_block_is_mineral == true { // caso in cui un blocco vicino sia un minerale e non sia il primo blocco generato
		var random int
		if position_x_previous_block == column-1 { // se il minerale è quello dietro allora la probabilità di generazione deve essere più bassa per impedire che si formino vene lunghe una fila dato che la generazione del mondo procede prima per le x
			random = rand.Intn(3)
		} else if position_y_previous_block == lenght-1 { // successivamente aumenta per fare in modo che se il blocco sia quello a fianco ci sia più probabilità di generare un minerale dello stesso tipo
			random = rand.Intn(5)
		} else { // qui indica se il minerale precedente si trova sotto il blocco da generare ed è più alta per fare si di avere vene che si sviluppano verso l'alto
			random = rand.Intn(7)
		}
		material, z, y, x, _, _ := mineral_recognition(high, lenght, column) // funzione per riconoscere il minerale e ne restituisce anche la posizione

		var previous_block block_position // struct con dentro materiale e posizione del minerale incontrato nelle vicinante del blocco da generare
		previous_block.material = material
		previous_block.pos_x = x
		previous_block.pos_y = y
		previous_block.pos_z = z

		if random == 0 { // calcola la probabilità di generare un minerale dello stesso tipo vicino
			current_block.material = stone
		} else {
			index_vein, is_in_vein := vein_recognition(previous_block) // dato in input il blocco precedente, verifica se si trova in una vena e restituisce un booleano e l'indice della vena in cui si trova (se non si trova in una vena l'indice sarà settato a zero zero zero zero, ma non è importante dato che viene inserito solo se il bool è vero)
			switch {
			case previous_block.material == coal && is_in_vein: // verifica se esite una vena di quel materiale,
				if total_veins[index_vein] < coal_vein { // verifica se nella vena c'è ancora spazio
					current_block.material = coal
					vein_update(current_block) // aggiorna l'indice della vena rendendo il blocco generato il blocco madre (per evitare che se il primissimo blocco della vena sia troppo lontano non venga più continuata la generazione)
				} else {
					current_block.material = stone
				}
			case previous_block.material == iron && is_in_vein:
				if total_veins[index_vein] < iron_vein {
					current_block.material = iron
					vein_update(current_block)
				} else {
					current_block.material = stone
				}
			case previous_block.material == gold && is_in_vein:
				if total_veins[index_vein] < gold_vein {
					current_block.material = gold
					vein_update(current_block)
				} else {
					current_block.material = stone
				}
			case previous_block.material == diamond && is_in_vein:
				if total_veins[index_vein] < diamond_vein {
					current_block.material = diamond
					vein_update(current_block)
				} else {
					current_block.material = stone
				}
			default:
				current_block.material = stone
			}
		}
	} else {
		current_block.material = stone
	}

	return current_block.material
}

func veins_nearby(z, y, x int) bool { // restituisce false se nelle vicinanze non ci sono minerali
	var empty_block block_position // genera un blocco vuoto che serve a verificare se nelle vicinanze ci sono vene minerarie
	empty_block.pos_x = x
	empty_block.pos_y = y
	empty_block.pos_z = z

	controller := zone_is_empty(empty_block) // verifica che la zona sia libera da minerali

	return controller
}

func mineral_recognition(high, lenght, column int) (previous_block, new_high, new_lenght, new_column int, not_first, mineral bool) {

	mineral = true
	not_first = true

	var controller_column bool = true
	var controller_lenght bool = true
	var controller_high bool = true

	if (column - 1) < 0 {
		controller_column = false
		not_first = false
	}
	if (lenght - 1) < 0 {
		controller_lenght = false
		not_first = false
	}
	if (high - 1) < 0 {
		controller_high = false
		not_first = false
	}

	if controller_column && controller_lenght && controller_high {

		if world[high][lenght][column-1] != stone && world[high][lenght][column-1] != bedrock && world[high][lenght][column-1] != dirt {
			previous_block = world[high][lenght][column-1]
			new_high = high
			new_lenght = lenght
			new_column = column - 1
		}

		if world[high][lenght-1][column] != stone && world[high][lenght-1][column] != bedrock && world[high][lenght-1][column] != dirt {
			previous_block = world[high][lenght-1][column]
			new_high = high
			new_lenght = lenght - 1
			new_column = column
		}

		if world[high-1][lenght][column] != stone && world[high-1][lenght][column] != bedrock && world[high-1][lenght][column] != dirt { // riconosce quali dei blocchi è il minerale e la sua posizione
			previous_block = world[high-1][lenght][column]
			new_high = high - 1
			new_lenght = lenght
			new_column = column
		}
		if previous_block == 0 {
			mineral = false
		}

	}
	return previous_block, new_high, new_lenght, new_column, not_first, mineral
}

func create_new_vein(block, x, y, z int) { // crea una nuova vena e utilizza come indice, il materiale e la posizione del primo minerale della vena generato
	var new_vein block_position // tipo struct
	new_vein.material = block
	new_vein.pos_x = x
	new_vein.pos_y = y
	new_vein.pos_z = z

	total_veins[new_vein] = 1
}

func vein_recognition(block block_position) (block_position, bool) { // riconosce a quale vena appartiene un minerale, restituisce l'indice della vena e il controllore che indica se è in una veno o no
	var controller bool
	var index block_position

	for i, _ := range total_veins { // scorre tutti gli indici
		var limit int = 5 // area di blocchi considerati vicini

		if i.material == block.material {

			difference_neg_x, _, _ := resising_limit(block, limit) // verifica che l'area da controllare non esca dalla mappa

			switch { // riconosce se il nuovo blocco è vicino al blocco madre nella map delle
			case i.pos_x >= difference_neg_x:
				controller = true
			case i.pos_y >= block.pos_y-1 && i.pos_y <= block.pos_y+1:
				controller = true
			case i.pos_z >= block.pos_z-1 && i.pos_z <= block.pos_z+1:
				controller = true
			}
			if controller {
				index = i
			}
		}
	}
	return index, controller
}

func zone_is_empty(block block_position) bool { // verifica che in un range dato non siano presenti minerali e restituisce vero solo se nel range è presente almeno un minerale
	var limit int = 5 // area di blocchi considerati vicini
	var controller bool
	difference_neg_x, difference_neg_y, difference_neg_z := resising_limit(block, limit) // verifica che l'area da controllare non esca dalla mappa

	for i := difference_neg_x; i < block.pos_x; i++ {
		for j := difference_neg_y; j < block.pos_y; j++ {
			for k := difference_neg_z + 1; k < block.pos_z; k++ {
				if world[k][j][i] != bedrock && world[k][j][i] != stone && world[k][j][i] != dirt {
					controller = true
					break
				}
			}
		}
	}
	return controller
}

func resising_limit(i block_position, limit int) (difference_neg_x, difference_neg_y, difference_neg_z int) { // impedisce che il range esca fuori dal mondo

	// le var difference pos, indicano il range x,y,z per capire se il blocco è vicino o meno

	if i.pos_x-limit < 0 { // serve a evitare che il range da controllare esca dai confini del mondo
		difference_neg_x = 0
	} else {
		difference_neg_x = i.pos_x - limit
	}

	if i.pos_y-limit < 0 {
		difference_neg_y = 0
	} else {
		difference_neg_y = i.pos_y - limit
	}

	if i.pos_z-limit < 0 {
		difference_neg_z = 0
	} else {
		difference_neg_z = i.pos_z - limit
	}

	return difference_neg_x, difference_neg_y, difference_neg_z
}

func vein_update(block block_position) { // sarà il nuovo indice nel caso si verificheranno le prossime condizioni

	i, controller := vein_recognition(block) // cerca a quale vena appartiene il minerale e ne restituisce l'indice nel caso controller sia vero

	if controller {
		old_value := total_veins[i]        // registra in una var temporanea il valore dell'indice attuale
		total_veins[block] = old_value + 1 // registra nella map delle vene minerarie le coordinate dell'ultimo blocco aggiunto, si somma uno per aumentare il valore della vena mineraria
		delete(total_veins, i)             // cancella il vecchio valore
	}
}

func data() { // stampa i dati del mondo
	var n_coal int
	var n_iron int
	var n_gold int
	var n_diamond int

	for i := 0; i < h_blocks; i++ {
		for j := 0; j < lenght_blocks; j++ {
			for k := 0; k < deep_blocks; k++ {
				if world[i][j][k] == coal {
					n_coal += 1
				} else if world[i][j][k] == iron {
					n_iron += 1
				} else if world[i][j][k] == gold {
					n_gold += 1
				} else if world[i][j][k] == diamond {
					n_diamond += 1
				}
			}
		}
	}
	fmt.Println("il numero di minerali di carbone è: ", n_coal)
	fmt.Println("il numero di minerali di ferro è: ", n_iron)
	fmt.Println("il numero di minerali di oro è: ", n_gold)
	fmt.Println("il numero di minerali di diamante è: ", n_diamond)
}

func print_levels() { // stampa il tipo di sezione desiderata
	var a string
	var mineral int
	var no_problem bool = true // controlla he non ci siano errori nel scegliere il minerale

	fmt.Println("scrivere in minuscolo e in inglese il minerale da vedere nella mappa (coal,iron,gold,diamond): ")
	fmt.Scan(&a)

	if a == "coal" {
		mineral = coal
	} else if a == "iron" {
		mineral = iron
	} else if a == "gold" {
		mineral = gold
	} else if a == "diamond" {
		mineral = diamond
	} else {
		fmt.Print("ERROR")
		no_problem = false
	}

	if no_problem {
		var what_level int
		var controller bool = true

		for controller {
			fmt.Println()
			fmt.Println()
			fmt.Println("inserisci:")
			fmt.Println("- 1 per stampare sezione parallela al piano")
			fmt.Println("- 2 sezione perpendicolare al piano e parallela a noi guardando da sinistra del mondo")
			fmt.Println("- 3 sezione perpendicolare al piano e perpendicolare a noi guardando da sinistra del mondo")
			fmt.Println("- qualsiasi altro numero per terminare")
			fmt.Println()
			fmt.Println()
			fmt.Println("legenda:")
			fmt.Println(" : blocco di aria")
			fmt.Println("@: blocco di roccia limite del mondo")
			fmt.Println("-: blocco di pietra")
			fmt.Println("0: blocco di terra")
			fmt.Println("#: blocco minerale scelto")

			fmt.Scan(&what_level)

			if what_level == 1 {
				for i := 0; i < h_blocks; i++ { // stampa sezione parallela al piano
					for j := 0; j < lenght_blocks; j++ {
						for k := 0; k < deep_blocks; k++ {
							if world[i][j][k] == bedrock {
								fmt.Print("@")
							} else if world[i][j][k] == mineral {
								fmt.Print("#")
							} else if world[i][j][k] == 0 {
								fmt.Print(" ")
							} else if world[i][j][k] == dirt {
								fmt.Print("0")
							} else {
								fmt.Print("-")
							}
						}
						fmt.Println()
					}
					fmt.Println()
					fmt.Println()
					fmt.Println("________________________________________________________________________________________________________________________")
					fmt.Println()
					fmt.Println()
				}
			} else if what_level == 2 {
				for k := 0; k < deep_blocks; k++ { // stampa sezione perpendicolare al piano e parallela a noi guardando da sinistra del mondo
					for i := h_blocks - 1; i >= 0; i-- {
						for j := 0; j < lenght_blocks; j++ {
							if world[i][j][k] == bedrock {
								fmt.Print("@")
							} else if world[i][j][k] == mineral {
								fmt.Print("#")
							} else if world[i][j][k] == 0 {
								fmt.Print(" ")
							} else if world[i][j][k] == dirt {
								fmt.Print("0")
							} else {
								fmt.Print("-")
							}
						}
						fmt.Println()
					}
					fmt.Println()
					fmt.Println()
					fmt.Println("________________________________________________________________________________________________________________________")
					fmt.Println()
					fmt.Println()
				}
			} else if what_level == 3 {
				for j := 0; j < lenght_blocks; j++ { // stampa sezione perpendicolare al piano e perpendicolare a noi guardando da sinistra del mondo
					for i := h_blocks - 1; i >= 0; i-- {
						for k := 0; k < deep_blocks; k++ {
							if world[i][j][k] == bedrock {
								fmt.Print("@")
							} else if world[i][j][k] == mineral {
								fmt.Print("#")
							} else if world[i][j][k] == 0 {
								fmt.Print(" ")
							} else if world[i][j][k] == dirt {
								fmt.Print("0")
							} else {
								fmt.Print("-")
							}
						}
						fmt.Println()
					}
					fmt.Println()
					fmt.Println()
					fmt.Println("________________________________________________________________________________________________________________________")
					fmt.Println()
					fmt.Println()
				}
			} else {
				controller = false
			}
		}

	}
}
