#!/usr/bin/env Rscript

library(ggplot2)
library(dplyr)
library(data.table)

main = function() {
	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=FALSE))
	colnames(data) = c("Chr", "Start", "End", "Pairing_Rate", "Cross", "CapH2")
	data$Cross = factor(data$Cross, levels=c("Pure D. mel", "Hybrid", "Hybrid Rescue"))
	levels(data$Cross) = c("Pure D. mel", "Hybrid", "Hybrid Rescue")

	data$CapH2 = factor(data$CapH2, levels=c("Peak", "Non-Peak"))
	levels(data$CapH2) = c("Peak", "Non-Peak")

	df.summary2 <- data %>%
		group_by(CapH2, Cross) %>%
		summarise(
			sd = sd(Pairing_Rate),
			Pairing_Rate = mean(Pairing_Rate)
		)

	g = ggplot(data, aes(CapH2, Pairing_Rate, fill = Cross), color = "black") +
		geom_col(data = df.summary2, position=position_dodge(.8), color = "black", width = .7) +
		theme_classic() +
		labs(x = "CapH2 Peak Type", y = "Pairing_Rate") +
		ggtitle("Pairing rate within or not within 1kb of CapH2 Peaks") +
		scale_fill_manual(values=c("#000000", "#888888", "#ffffff"), breaks = c("Pure D. mel", "Hybrid", "Hybrid Rescue"))

	pdf(args[2], width=4, height=3)
		print(g)
	dev.off()

	g2 = g + geom_errorbar(data = df.summary2, aes(ymin = Pairing_Rate - sd, ymax = Pairing_Rate + sd), position = position_dodge(.8), color = "black", width = .2)

	pdf(args[3], width=4, height=3)
		print(g2)
	dev.off()
}

main()

# #!/usr/bin/env Rscript
# 
# library(ggplot2)
# 
# main = function() {
# 	data = read.table("count_avgs_fmt.txt", header=FALSE)
# 	colnames(data) = c("Breed", "Number_Shared", "Reps_Counted", "Counts", "Average_Per_Comparison")
# 	data$Breed = factor(data$Breed, levels = c("Black", "White", "Runt", "Figurita"))
# 	levels(data$Breed) = c("White", "Black", "Runt", "Figurita")
# 	data = data[data$Number_Shared > 0,]
# 
# 	g = ggplot(data, aes(Number_Shared, Average_Per_Comparison, fill=Breed), color = "black") + 
# 		geom_col(position="dodge", color = "black") +
# 		theme_classic() +
# 		labs(x = "Replicates", y = "Average shared peaks") +
# 		ggtitle("Counts of shared pFst peaks") +
# 		scale_fill_manual(values=c("#000000", "#ffffff", "#555555", "#aaaaaa"), breaks = c("Black", "White", "Runt", "Figurita"))
# 
# 	pdf("counts_avg_final.pdf", width=4, height=3)
# 	print(g)
# 	dev.off()
# }
# 
# main()


# #!/usr/bin/env Rscript
# 
# library(ggplot2)
# library(data.table)
# 
# main = function() {
# 	args = commandArgs(trailingOnly=TRUE)
# 	data = as.data.frame(fread(args[1], header=TRUE))
# 
# 	data$Cross = factor(data$Cross, levels=c("Pure D. mel", "Hybrid"))
# 
# 	p = ggplot(data = data[data$snpclass != "bitted_unsel",], aes(snpclass, Slope)) +
# 		geom_boxplot() +
# 		labs(x = "Site type", y = "Pairing rate") +
# 		scale_x_discrete(name = "SNP class", breaks = breaks, labels = labels) +
# 		theme_bw() +
# 		theme(text = element_text(size=18))
# 
# 	pdf(args[2], height = 3, width = 4)
# 		print(p)
# 	dev.off()
# 
# 	# breaks = c("unbitted", "bitted", "unbitted_unsel")
# 	# labels = c("Unbitted", "Bitted", "Unselected")
# 	# data = data[data$snpclass != "bitted_unsel",]
# 		# lims(y=c(-0.012, 0.018)) +
# }
# 
# main()
# 


# #!/usr/bin/env Rscript
# 
# library(ggplot2)
# library(dplyr)
# 
# main = function() {
# 	data = read.table("count_raws_fmt.txt", header=FALSE)
# 	colnames(data) = c("Breed", "Number_Shared", "Counts")
# 	data$Breed = factor(data$Breed, levels = c("White", "Black", "Runt", "Figurita"))
# 	levels(data$Breed) = c("White", "Black", "Runt", "Figurita")
# 	data = data[data$Number_Shared > 0,]
# 
# 	df.summary2 <- data %>%
# 	  group_by(Number_Shared, Breed) %>%
# 	  summarise(
# 	    sd = sd(Counts),
# 	    Counts = mean(Counts)
# 	  )
# 	df.summary2
# 
# 	g = ggplot(data, aes(Number_Shared, Counts, fill=Breed), color = "black") + 
# 		geom_col(data = df.summary2, position=position_dodge(.8), color = "black", width = .7) +
# 		geom_errorbar(data = df.summary2, aes(ymin = Counts - sd, ymax = Counts + sd), position = position_dodge(.8), color = "black", width = .2) +
# 		theme_classic() +
# 		labs(x = "Replicates", y = "Average shared peaks") +
# 		ggtitle("Counts of shared pFst peaks") +
# 		scale_fill_manual(values=c("#ffffff", "#444444", "#888888", "#cccccc"), breaks = c("White", "Black", "Runt", "Figurita"))
# 
# 	pdf("counts_avg_final2.pdf", width=4, height=3)
# 	print(g)
# 	dev.off()
# }
# 
# main()
# 
